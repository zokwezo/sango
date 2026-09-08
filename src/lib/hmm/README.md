## Hidden Markov Model (HMM) with Kneser-Ney Smoothing

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

### Advantages of an HMM over CANINE:

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

### Disadvantages and mitigation

Although the HMM architecture is robust for a small training corpus, it suffers from overreliance on tag transitions over emission probabilities and using relies on finicky tuning of damping hyperparameters and favors common tokens and disfavors rare tokens and especially hapax legomena.

### Mitigation

Using Kneser-Ney smoothing on the HMM emission matrix solves the problem of the transition matrix overriding the emissions without the need for exponential dampening and the required careful tuning of hyperparameters:

1. **Protects the Model from Idiom Over-weighting:** If the small training corpus contains a highly specific phrase (e.g., a specific noun always paired with a specific technical tag), Kneser-Ney prevents that noun from bleeding out and dominating other tags when encountered in a new sentence structure.
2. **Smart Allocation to Unseen Tokens:** For out-of-vocabulary or rare tokens, the model will fallback to a continuation probability. It will ask: "Is this tag/label generally a versatile one that accepts many kinds of words?" If a hidden state is highly versatile (like a generic NOUN or VERB state), it gets a higher share of the smoothed probability mass than a highly restrictive state (like a specific punctuation state).

### Algorithm

Kneser-Ney uses **Absolute Discounting** as its base but replaces raw counts with continuation probabilities for lower-order histories.

For a bigram/emission context, the probability of a word w given a state/context c is calculated as:

$$P_{KN}(w\mid c)=\frac{\max (\text{count}(c,w)-d,0)}{\text{count}(c)}+\lambda (c)\cdot P_{continuation}(w)$$

1. **The Left Term (Absolute Discounting):** It takes the raw count of the pair, subtracts a flat discount d (usually between 0.5 and 0.75), and divides by the total context count.
2. **The Right Term (λ Back-off Weight):** This is the normalized weight of the "stolen" probability mass we accumulated by subtracting d.
3. **The Continuation Probability ($P_{continuation}$):** Instead of distributing the stolen mass uniformly (like Absolute Discounting does), Kneser-Ney distributes it based on how many unique contexts the word w has appeared in:

$$P_{continuation}(w)=\frac{\text{Number\ of\ unique\ contexts\ }c\text{\ where\ }(c,w)\text{\ occurs}}{\text{Total\ number\ of\ unique\ bigram\ types\ in\ the\ corpus}}$$

## Quality Metrics

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

