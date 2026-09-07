# Hidden Markov Model (HMM)

This directory contains code to perform **token classification (sequence labeling)**.

The `demo` script contains a pipeline to train, test, and evaluate pitch diacritic restoration.

## Model design

In the current model design, the tokens are Sango words (after pitch diacritics are removed) and the labels (also known as tags) are the word pitch contour which is the sequence of pitches (low, medium, high), one for each vowel in the word.

This is not the only possible choice. For instance, a token might instead be a Sango syllable which would have several obvious advantages:
* They have a unique label (a single pitch for the single vowel in a syllable), which allows no possibility of a pitch sequence with a length different from the number of syllables in a word.
* The emission matrix is a lot smaller, which would give a less noisy model given the small size of the training corpus.

The disadvantages are likely much stronger:
* Syllables yield no natural emission matrix, inviting the model to overtrain and hallucinate.
* Individual pitch transitions are relatively uninformative, whereas word pitch contours of common words are a stronger predictor.

The final decision on which tokens and labels are most predictive is a work in project and essentially an empirical decision.

## Model architecture

A Hidden Markov Model (HMM) is used rather than an encoder model (such as the Google [CANINE](https://huggingface.co/docs/transformers/en/model_doc/canine) model, much less an encoder/decoder model (such as a sequence-to-sequence model) due to the extremely low size of the available Sango training corpus with reliable diacritics, since almost the entire extant corpus is either lacking in diacritics or has a high error rate.

The advantages of a HMM over CANINE are:

* Resistance to Severe Overfitting (Parameter Efficiency)
  - An HMM uses a relatively small number of parameters—specifically, its transition, emission, and initial state probability matrices. Because it has low model capacity, it can be mathematically optimized using small datasets via the Baum-Welch algorithm without memorizing the noise. [1](https://huggingface.co/blog/RDTvlokip/from-scratch-vs-pre-trained)
  - With 121M parameters, CANINE has vast representation capacity. When exposed to a tiny corpus, it suffers from severe overfitting. It will simply memorize the small dataset rather than learning the foundational syntax, phonetics, or morphology of the Sango language. [2](https://huggingface.co/transformers/v4.12.5/model_doc/canine.html), [3](https://huggingface.co/blog/RDTvlokip/from-scratch-vs-pre-trained)
* Elimination of Pre-Training Constraints
  - An HMM is trained entirely from scratch on the target corpus. It relies strictly on the mathematical statistics of the sequences you provide. [4](https://jonathan-hui.medium.com/speech-recognition-gmm-hmm-8bb5eff8b196)
  - CANINE relies heavily on cross-lingual transfer learning from its 104 pre-trained languages. Since the low-resource Sango language uses diacritics to represent pitch accents highly specific to each token, CANINE's pre-trained vocabulary-free layers cannot successfully map it. Attempting to train CANINE from scratch on the small corpus available with accurate diacritics, it will fail completely, as Transformers generally require millions or billions of tokens to converge. [5](https://www.researchgate.net/publication/393899105_TRANSFORMER_MODELS_IN_LOW-RESOURCE_LANGUAGE_TRANSLATION), [6](https://www.sciencedirect.com/science/article/pii/S2949719125000330), [7](https://medium.com/@deependraiimb/training-a-large-language-model-on-your-own-data-3b3e4398e400), [8](https://huggingface.co/google/canine-c), [9](https://aclanthology.org/2023.conll-1.35/)
* Better Handling of Character Sequence Lengths
  - HMMs process linear, step-by-step state transitions between tokens efficiently over local context window constraints (like a standard n-gram or bigram approach). [10](https://ar5iv.labs.arxiv.org/html/2112.10508)
  - Because CANINE is tokenization-free and reads text character-by-character, its sequence lengths are exceptionally long. To combat this, it uses strided convolutions to downsample characters. On a tiny corpus, the model cannot effectively learn how to stitch these downsampled character sequences back into semantic meaning. [11](https://medium.com/@amey.cmyo/how-googles-canine-2021-model-learns-languages-without-even-splitting-words-dd456e38e66f), [12](https://aclanthology.org/2022.tacl-1.5.pdf)

*The above section is taken from Google Gemini, see the original thread [here](https://share.google/aimode/pcOyHPe5W4VRnRghe).*

## Implementation

In an HMM, Transition and Emission probabilities are the core mathematical rules used to determine which label belongs to which token.

> *TODO: When a larger training corpus becomes available, switch to using tokens instead (after separating prefixes and suffixes into their own tokens).*

Viterbi Algorithm Decoding is used for a globally optimal sequence of labels rather than making greedy token-by-token guesses.

Laplace smoothing is used for tokens it encounters during prediction but didn't see during training, though this will occur only for foreign or misspelled tokens.

### Hidden model parameters

The transition, emission, and starting probabilities are estimated from training sentences.

#### 1. (State) Transition Probabilities: P(Label₂ | Label₁)

Transition probability is the likelihood of moving from one specific label to another label. It answers the question: Given the current label, what is the probability of the next label?

* What it measures: The grammar, structure, and syntax patterns of the language.
* How it is calculated:
$$\text{Transition Probability} = \frac{\text{Number of times Label₁ is followed by Label₂}}{\text{Total occurrences of Label₁}}$$

#### 2. Emission (or Posterior) Probabilities: P(Token | Label)

Emission probability is the likelihood that a specific hidden label will emit (or produce) a specific visible token. It answers the question: Given a specific label, how likely is it to apply to this exact token?

* What it measures: The vocabulary and semantic meaning associated with a label.
* How it is calculated:
$$\text{Emission Probability} = \frac{\text{Number of times a specific Token is labeled with a specific Label}}{\text{Total occurrences of that Label}}$$ 

#### 3. Starting (or Marginal) Probabilities: P(Label₀)

Starting probability is the likelihood that a specific hidden label will start a sentence. This grounds the backwards recursion and answers the question: Given a specific label, how likely is it to start the sentence?

* What it measures: The vocabulary and semantic meaning associated with a label.
* How it is calculated:
$$\text{Starting Probability} = \frac{\text{Number of times a specific Label starts a sentence}}{\text{Number of training sentences}}$$ 

#### How They Work Together (The Viterbi Decoding)

The above visible emission and starting parameters (also known as the posterior and marginal probabilities, respectively) are unfortunately the inverse of the (hidden to us) conditional and prior probabililties that we actually need to predict the next label, and the latter are reconstructed via Bayes' Theorem by finding the maximum likelihood path backwards through a transition chain of labels from the end of a sentence to its start. This is the essence of the Viterbi implementation of the Hidden Markov Model.

By calculating all possible path combinations using these two probability scores, the algorithm determines the most accurate sequence labels for the text.

### Compensating for model sampling error

We do not have access to the true population statistics of Sango sentences, and instead have to make do with the noisily estimated sample statistics obtained from training on a small corpus.

To improve the model stability due to sampling error, Laplace smoothing and log-probability are used:

#### Laplace Smoothing

If an unknown token shows up during evaluation, standard calculation would multiply by 0% probability, wiping out the entire sequence score.

Instead, the training code adds 1 to every event counter and increases the normalization denominator by the total vocabulary or label pool size. Unseen tokens fallback safely onto a tiny, non-zero probability instead of failing outright.

#### Log-Probability Arithmetic

Multiplying numbers under 1.0 repeatedly across long sentences creates extremely small fractions (e.g. $0.1^{20} = 10^{-20}$), causing system floating-point numbers to rounds down to absolute $0.0$ (Underflow).

All probabilities are transformed using math.Log(). The math rules change from multiplication to simple addition: $\log (A\times B)=\log (A)+\log (B)$. This ensures mathematical stability even across sentences containing thousands of tokens.

### Quality Metrics

The model can be applied to test sentences not included in the training corpus to generate the following quality metrics:
* Precision (Positive Predictive Value)
  - fraction of predicted that are correct
  - High precision means the model rarely flags standard tokens incorrectly.
* Recall (Sensitivity)
  - fraction of correct that are predicted
  - High recall means the model rarely misses target entities.
* F1-Score
  - The harmonic mean balancing Precision and Recall into a single quality score.
  - It is the gold standard metric for optimizing sequence classifiers.

*NOTE: If a test sentence were in the training corpus, the quality metrics would be artificially inflated since the model would have essentially memorized the sentence.*
