from preprocess import preprocess_image
from utils import untokenize, checkpoint
from model.models import Encoder, Decoder
import torch
import numpy as np
from PIL import Image
import io
from torchvision.models.resnet import ResNet50_Weights

device = torch.device("cpu")

# todo: docstring
def caption_image(encoder, decoder, image, image_processor, word_map):
    encoder.eval()
    decoder.eval()
    
    # preprocess image
    img = preprocess_image(image, image_processor)

    # convert to PyTorch tensor and add batch dimension
    img = torch.tensor(np.expand_dims(img, axis=0)).float().to(device)

    # extract features
    features = encoder(img)

    # generate caption
    vocab_size = len(word_map)
    caption, _ = decoder.caption(features, start_token=vocab_size-2, end_token=vocab_size-1)

    
    return untokenize.untokenize_caption(caption, word_map)

# TODO: THIS IS A TEMP FUNCTION
def load_image_as_byte_array(image_path):
    # Open the image file
    with Image.open(image_path) as img:
        # Convert image to byte array
        img_byte_array = io.BytesIO()
        img.save(img_byte_array, format=img.format)
        return img_byte_array.getvalue()

if __name__ == "__main__":
    device = torch.device("cpu")

    # Load word map
    word_map = checkpoint.load_json("./wordmap/wordmap.json")

    # Load the models
    decoder = Decoder(vocab_size=len(word_map), attention_dim=512, embedding_dim=2048, decoder_dim=512, num_layers=2).to(device)
    encoder = Encoder().to(device)
    checkpoint.load_checkpoint(encoder, decoder, checkpoint_path='/home/gio/image_classifier/models/checkpoint.pth')

    # Initialize image transform
    weights = ResNet50_Weights.DEFAULT
    image_processor = weights.transforms()

    # load image
    image_byte_array = load_image_as_byte_array("./test_pants.png")

    # caption image
    caption = caption_image(encoder, decoder, image_byte_array, image_processor, word_map)

    print(caption)