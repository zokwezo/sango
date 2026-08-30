# Hidden Markov Model (HMM)

This directory contains code produced by Google Gemini to perform **token classification (sequence labeling)**.

In sequence labeling models like a Hidden Markov Model (HMM), Transition and Emission probabilities are the core mathematical rules used to determine which tag belongs to which word.

To understand them, think of the tags (like PERSON, LOCATION, VERB) as hidden states you cannot see directly, and the words themselves as the observed data.

The advantages of a Hidden Markov Model are:

* Custom Labeled Data Ready: You feed custom sequences to the Train() function.
* Viterbi Algorithm Decoding: Computes the globally optimal sequence of labels rather than making greedy word-by-word guesses.
* Basic Smoothing: Allocates a default low probability (1e-5) for words it encounters during prediction but didn't see during training.

------------------------------

## 1. Transition Probabilities: P(Tag₂ | Tag₁)

Transition probability is the likelihood of moving from one specific tag to another tag. It answers the question: "Given the current tag, what is the probability of the next tag?"

* What it measures: The grammar, structure, and syntax patterns of a language.
* Why it matters: In English text, a person's first name is often followed by a person's last name, and a determiner (like "the") is often followed by a noun. It is highly unlikely for a location tag to immediately follow a first name tag without a verb or preposition in between.
* How it is calculated:
$$\text{Transition Probability} = \frac{\text{Number of times Tag₁ is followed by Tag₂}}{\text{Total occurrences of Tag₁}}$$ 
* Example: If the tag B-PER (Begin Person) appears 100 times in your training data, and it is followed by I-PER (Inside Person) 80 times, the transition probability $P(\text{I-PER} \mid \text{B-PER})$ is 0.80 (80%).

------------------------------

## 2. Emission Probabilities: P(Word | Tag)
Emission probability is the likelihood that a specific hidden tag will "emit" or produce a specific visible word. It answers the question: "Given a specific tag, how likely is it to look like this exact word?"

* What it measures: The vocabulary and semantic meaning associated with a tag.
* Why it matters: The tag B-LOC (Begin Location) is highly likely to emit words like "Paris", "London", or "Cupertino", but almost never emits words like "ate", "running", or "apple".
* How it is calculated:
$$\text{Emission Probability} = \frac{\text{Number of times a specific Word is labeled with a specific Tag}}{\text{Total occurrences of that Tag}}$$ 
* Example: If the tag B-LOC appears 50 times in your dataset, and 10 of those times the corresponding word is "Paris", the emission probability $P(\text{"Paris"} \mid \text{B-LOC})$ is 0.20 (20%).

------------------------------

## How They Work Together (The Viterbi Decoding)
When you pass a sentence like "Steve visited Cupertino" to the model, it multiplies these two types of probabilities together across the entire sequence to find the chain of tags that yields the highest overall score:

   1. Initial Probability: How likely is it that the sentence starts with B-PER?
   2. Emission: How likely is B-PER to produce the word "Steve"?
   3. Transition: How likely is it to transition from B-PER to O (Outside/Verb)?
   4. Emission: How likely is O to produce the word "visited"?

By calculating all possible path combinations using these two probability scores, the algorithm determines the most accurate sequence labels for your text.
If you are interested, I can show you how to:

* Implement Laplace smoothing in Go so unknown words don't break your emission scores.
* Print out the raw probability matrix from the previous Go code to see these numbers in action.
* Convert these calculations into Log-probabilities to prevent floating-point underflow errors on long sentences.

## Upgrades

### Laplace (Add-1) Smoothing

#### The Problem:

If an unknown word like "London" shows up during evaluation, standard calculation multiplies by 0% probability, wiping out the entire sequence score.

#### The Fix:

The training code adds 1.0 to every event counter and increases the normalization denominator by the total vocabulary or tag pool size. Unseen words fallback safely onto a tiny, non-zero probability ("<UNKNOWN>") instead of failing.

### Log-Probability Arithmetic

#### The Problem:

Multiplying numbers under 1.0 repeatedly across long sentences creates extremely small fractions (e.g., \(0.1^{20} = 10^{-20}\)), causing system floating-point numbers to rounds down to absolute 0.0 (Underflow).

#### The Fix:

All probabilities are transformed using math.Log(). The math rules change from multiplication to simple addition:\(\log (A\times B)=\log (A)+\log (B)\)This ensures mathematical stability even across sentences containing thousands of words.

### Matrix Inspection (PrintMatrices)

Transformed log math values are reverted back to standard percentage form using math.Exp().This allows you to visually audit your grammar rules (Transition) and dictionary rules (Emission) directly inside the console terminal.

### Evaluation metrics

* Precision (Positive Predictive Value)
  - fraction of predicted that are correct
	- High precision means the model rarely flags standard words incorrectly.
* Recall (Sensitivity)
  - fraction of correct that are predicted
	- High recall means your model rarely misses target entities.
* F1-Score
  - The harmonic mean balancing Precision and Recall into a single quality score.
	- It is the gold standard metric for optimizing sequence classifiers.
