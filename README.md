# image_captioner_services

Captioning flow:

a user uploads an image (frontend sends image to endpoint managed by Amazon API Gateway)

API gateway triggers a lambda function that saves image to S3 bucket

S3 event triggers lambda function that sends image to SQS with image details

SQS message triggers lambda function that processes image with Docker container, this caption is then sent to S3 bucket
