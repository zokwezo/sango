## Hidden Markov Model (HMM) with Kneser-Ney Smoothing

### Advantages

Using Kneser-Ney smoothing on the HMM emission matrix solves the problem of the transition matrix overriding the emissions without the need for careful tuning of hyperparameters:

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
