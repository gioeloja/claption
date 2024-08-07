import torch
import numpy as np
from utils import preprocess, untokenize

device = torch.device("cpu")

# todo: docstring
def caption_image(encoder, decoder, image, image_processor, word_map):
    encoder.eval()
    decoder.eval()
    
    # preprocess image
    img = preprocess.preprocess_image(image, image_processor)

    # convert to PyTorch tensor and add batch dimension
    img = torch.tensor(np.expand_dims(img, axis=0)).float().to(device)

    # extract features
    features = encoder(img)

    # generate caption
    vocab_size = len(word_map)
    caption, _ = decoder.caption(features, start_token=vocab_size-2, end_token=vocab_size-1)

    
    return untokenize.untokenize_caption(caption, word_map)