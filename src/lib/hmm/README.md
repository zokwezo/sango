# Hidden Markov Model (HMM)

This directory contains code to perform **token classification (sequence labeling)**, where the tokens are Sango syllables (after pitch diacritics are removed) and the labels are pitch level (low, medium, high).

The `demo` script contains a pipeline to train, test, and evaluate pitch diacritic restoration.

## Model architecture

A Hidden Markov Model (HMM) is used rather than an encoder model (such as the Google [CANINE](https://huggingface.co/docs/transformers/en/model_doc/canine) model, much less an encoder/decoder model (such as a sequence-to-sequence model) due to the extremely low size of the available Sango training corpus with reliable diacritics, since almost the entire extant corpus is either lacking in diacritics or has a high error rate.

The advantages of a HMM over CANINE are:

* Resistance to Severe Overfitting (Parameter Efficiency)
  - An HMM uses a relatively small number of parameters—specifically, its transition, emission, and initial state probability matrices. Because it has low model capacity, it can be mathematically optimized using small datasets via the Baum-Welch algorithm without memorizing the noise. [1](https://huggingface.co/blog/RDTvlokip/from-scratch-vs-pre-trained)
  - With 121M parameters, CANINE has vast representation capacity. When exposed to a tiny corpus, it suffers from severe overfitting. It will simply memorize the small dataset rather than learning the foundational syntax, phonetics, or morphology of the Sango language. [2](https://huggingface.co/transformers/v4.12.5/model_doc/canine.html), [3](https://huggingface.co/blog/RDTvlokip/from-scratch-vs-pre-trained)
* Elimination of Pre-Training Constraints
  - An HMM is trained entirely from scratch on the target corpus. It relies strictly on the mathematical statistics of the sequences you provide. [4](https://jonathan-hui.medium.com/speech-recognition-gmm-hmm-8bb5eff8b196)
  - CANINE relies heavily on cross-lingual transfer learning from its 104 pre-trained languages. Since the low-resource Sango language uses diacritics to represent a pitch accent highly specific to each syllable, CANINE's pre-trained vocabulary-free layers cannot successfully map it. Attempting to train CANINE from scratch on the small corpus available with accurate diacritics, it will fail completely, as Transformers generally require millions or billions of syllables to converge. [5](https://www.researchgate.net/publication/393899105_TRANSFORMER_MODELS_IN_LOW-RESOURCE_LANGUAGE_TRANSLATION), [6](https://www.sciencedirect.com/science/article/pii/S2949719125000330), [7](https://medium.com/@deependraiimb/training-a-large-language-model-on-your-own-data-3b3e4398e400), [8](https://huggingface.co/google/canine-c), [9](https://aclanthology.org/2023.conll-1.35/)
* Better Handling of Character Sequence Lengths
  - HMMs process linear, step-by-step state transitions between syllables efficiently over local context window constraints (like a standard n-gram or bigram approach). [10](https://ar5iv.labs.arxiv.org/html/2112.10508)
  - Because CANINE is tokenization-free and reads text character-by-character, its sequence lengths are exceptionally long. To combat this, it uses strided convolutions to downsample characters. On a tiny corpus, the model cannot effectively learn how to stitch these downsampled character sequences back into semantic meaning. [11](https://medium.com/@amey.cmyo/how-googles-canine-2021-model-learns-languages-without-even-splitting-words-dd456e38e66f), [12](https://aclanthology.org/2022.tacl-1.5.pdf)

*The above section is taken from Google Gemini, see the original thread [here](https://share.google/aimode/pcOyHPe5W4VRnRghe).*

## Implementation

In an HMM, Transition and Emission probabilities are the core mathematical rules used to determine which pitch belongs to which syllable. Syllables are used as tokens rather than whole words to reduce the dimensionality due to the small size of the training corpus.

> *TODO: When a larger training corpus becomes available, switch to using words instead (after separating prefixes and suffixes into their own tokens).*

Viterbi Algorithm Decoding is used for a globally optimal sequence of labels rather than making greedy syllable-by-syllable guesses.

Laplace smoothing is used for syllables it encounters during prediction but didn't see during training, though this will occur only for foreign or misspelled words.

### Transition and Emission

#### 1. Transition Probabilities: P(Pitch₂ | Pitch₁)

Transition probability is the likelihood of moving from one specific pitch to another pitch. It answers the question: Given the current pitch, what is the probability of the next pitch?

* What it measures: The grammar, structure, and syntax patterns of the language.
* Why it matters: In Sango text, certain syllables tend to frequently follow or precede specific other syllables (especially common monosyllabic words or prefixes). There is also vowel pitch harmony with certain parts of speech, for instance a syllable preceding the gerund suffix -ngɔ̈ (written -ngö in the standard orthography) is extremely likely to be mid-pitch.
* How it is calculated:
$$\text{Transition Probability} = \frac{\text{Number of times Pitch₁ is followed by Pitch₂}}{\text{Total occurrences of Pitch₁}}$$
* Example: If the pitch $\text{LO}$ (low pitch) appears 100 times in the training data, and it is followed by $\text{HI}$ (high pitch) 80 times, the transition probability $P(\text{HI} \mid \text{LO})$ is 0.80 (or 80%).

#### 2. Emission Probabilities: P(Syllable | Pitch)

Emission probability is the likelihood that a specific hidden pitch will emit (or produce) a specific visible syllable. It answers the question: Given a specific pitch, how likely is it to apply to this exact syllable?

* What it measures: The vocabulary and semantic meaning associated with a pitch.
* Why it matters: The pitch $\text{MID}$ (for mid pitch) is unlikely to emit syllables in highly frequent high-pitched monosyllabic words such as **tî**, **nî**, and **sô**.
* How it is calculated:
$$\text{Emission Probability} = \frac{\text{Number of times a specific Syllable is labeled with a specific Pitch}}{\text{Total occurrences of that Pitch}}$$ 
* Example: If a MID pitch appears 50 times in the dataset, and 10 of those times the corresponding syllable is **ngo**, the emission probability $P(\text{ngo} \mid \text{MID})$ is 0.20 (20%).

#### How They Work Together (The Viterbi Decoding)

When you pass a sentence like **na mbâgë tî wâlï** to the model, it multiplies these two types of probabilities together across the entire sequence to find the chain of pitchs that yields the highest overall score:

   1. Initial Probability: How likely is it that the sentence starts with low pitch?
   2. Emission: How likely is a low pitch to produce the syllable **na**?
   3. Transition: How likely is it to transition from low pitch to high pitch?
   4. Emission: How likely is a high pitch to produce the pitch-agnostic syllable **mba**?

By calculating all possible path combinations using these two probability scores, the algorithm determines the most accurate sequence labels for the text.

### Laplace Smoothing

If an unknown syllable like **fri** were to show up during evaluation (such as in the non-standard **Afrika** instead of the standard **Afirika**), standard calculation would multiply by 0% probability, wiping out the entire sequence score.

Instead, the training code adds 1 to every event counter and increases the normalization denominator by the total vocabulary or pitch pool size. Unseen syllables fallback safely onto a tiny, non-zero probability (indicated with an underdot) instead of failing outright.

### Log-Probability Arithmetic

Multiplying numbers under 1.0 repeatedly across long sentences creates extremely small fractions (e.g. $0.1^{20} = 10^{-20}$), causing system floating-point numbers to rounds down to absolute $0.0$ (Underflow).

All probabilities are transformed using math.Log(). The math rules change from multiplication to simple addition: $\log (A\times B)=\log (A)+\log (B)$. This ensures mathematical stability even across sentences containing thousands of syllables.

### Quality Metrics

The model can be applied to test sentences not included in the training corpus to generate the following quality metrics:
* Precision (Positive Predictive Value)
  - fraction of predicted that are correct
  - High precision means the model rarely flags standard syllables incorrectly.
* Recall (Sensitivity)
  - fraction of correct that are predicted
  - High recall means the model rarely misses target entities.
* F1-Score
  - The harmonic mean balancing Precision and Recall into a single quality score.
  - It is the gold standard metric for optimizing sequence classifiers.

*NOTE: If a test sentence were in the training corpus, the quality metrics would be artificially inflated since the model would have essentially memorized the sentence.*
