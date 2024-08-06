"""
    Returns the word corresponding to a given index

    Parameters:
        idx: idx to map
        wordmap: dictionary where words (keys) map to indices (values)
"""
def idx_to_word(idx, wordmap):
    for word, word_idx in wordmap.items():
        if word_idx == idx:
            return word
    return None

# TODO: docstring
def untokenize_caption(tokens, word_map):
    untokenized_caption = ""

    for i in tokens:
        pred_word = idx_to_word(i, word_map)

        if pred_word and pred_word != "<start>" and pred_word != "<end>" and pred_word != "<pad>":
            untokenized_caption += pred_word + " "
    
    return untokenized_caption
