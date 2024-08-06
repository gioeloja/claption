from flask import Flask, request, jsonify
import io

app = Flask(__name__)

@app.route('/caption-image', methods=['POST'])
def handler():
    if 'file' not in request.files:
        return jsonify({"error": "No file part in the request"}), 400

    file = request.files['file']

    if file.filename == '':
        return jsonify({"error": "No file selected for uploading"}), 400

    try:
        image_data = file.read()  # Read the image data as bytes

        # Here you can process the image data (image_data) as needed
        # For now, we'll just return a success message

        return jsonify({"message": "Image received", "image_data_length": len(image_data)}), 200
    except Exception as e:
        return jsonify({"error": str(e)}), 500

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=80)
