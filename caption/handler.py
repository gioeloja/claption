from flask import Flask, request, jsonify
import torch
from model.models import Encoder, Decoder
from utils import checkpoint, caption_image

from torchvision.models.resnet import ResNet50_Weights

app = Flask(__name__)

device = torch.device("cpu")

# word map used for creating captions from a select list of words
word_map = checkpoint.load_json("./wordmap/wordmap.json")

# initialize models and load weights from checkpoint
decoder = Decoder(vocab_size=len(word_map), attention_dim=512, embedding_dim=2048, decoder_dim=512, num_layers=2).to(device)
encoder = Encoder().to(device)
checkpoint.load_checkpoint(encoder, decoder, checkpoint_path='./model/checkpoint.pth')

# image preprocessor
weights = ResNet50_Weights.DEFAULT
image_processor = weights.transforms()

@app.route('/caption-image', methods=['POST'])
def handler():
    if 'file' not in request.files:
        return jsonify({"error": "No file included in request"}), 400

    file = request.files['file']

    if file.filename == '':
        return jsonify({"error": "No file uploaded"}), 400

    try:
        image_data = file.read()  # image data as bytes

        caption = caption_image.caption_image(encoder, decoder, image_data, image_processor, word_map)

        return jsonify({"message": caption})
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
