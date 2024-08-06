from PIL import Image
import io

"""
Process image by running it through a transform object.

Parameters:
    - image: image in []byte form
    - transform: image transformer

returns: transformed PIL image
"""
def preprocess_image(image: bytes, transform):
    # convert to a PIL Image
    img = Image.open(io.BytesIO(image))

    # todo: try catch
    return transform(img)