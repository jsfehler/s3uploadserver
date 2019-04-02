# s3uploadserver
Web server that uploads files to Amazon S3

The server assumes the environment it's being run on has AWS API access.

## Usage

GET: Will always return 200, plus a message about the server.

POST: Will upload multipart form-data to Amazon S3. Returns a list of files that were uploaded.

Expects a filename, which is used for the upload.

### Metadata

If a form called "metadata" is included, that data will be uploaded as a separate file named with the filename plus _metadata.log, located next to the uploaded file.

eg: `images/myfile.png` and `images/myfile_metadata.log`

#### S3Root parameter

If the metadata contains a parameter called `S3Root` with a string, then this will be added as a prefix to the upload location.

eg: If the file path is `images/myfile.png` and S3Root is `001` then the file will be uploaded to `001/images/myfile.png`

### Running via Docker

`docker run --rm -it -p 8081:8081 -e AWS_ACCESS_KEY_ID={your_access_key_id} -e AWS_SECRET_ACCESS_KEY={your_secret_access_key} jsfehler/s3uploadserver -bucket={your_s3_bucket}`

### Flags

- port: Port to launch server on. Default is `8081`.
- bucket: S3 bucket name to send files to.
- region: AWS region. Default is `us-east-1`.
