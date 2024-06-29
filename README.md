# image_captioner_services


Development Plan:
1. Make sure the captioner works locally. Get it setup so that it's only using CPU and see how long image captioning takes.
2. Containerize the model using Docker. Figure out how to feed images into this Docker image.
3. Mock AWS locally.
   - Upload images to a local S3 bucket. This will trigger a lambda function that sends a message to SQS.
   - Lambda function will look at SQS for messages and then fetch images from the local S3.
   - Docker image is triggered to process the image and store the captions back in S3.
4. Figure out what else I can do before actually doing AWS configuration. Some ideas...
   - Create a frontend where users can upload images
   - Maybe some functionality that stores a user's captions
       - This will require accounts (auth)
