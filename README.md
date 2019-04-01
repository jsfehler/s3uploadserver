# s3uploadserver
Web server that uploads files to Amazon S3

The server assumes the machine it's being run on has AWS API access.

## Usage

GET: Will always return 200, plus a message about the server.

POST: Will upload multipart form-data to Amazon S3. Returns a list of files that were uploaded.

Expects a filename, which is used for the upload.

If a form called "metadata" is included, that data will be uploaded as a separate file named with the filename plus _metadata.log, located next to the uploaded file.

eg: `images/myfile.png` and `images/myfile_metadata.log`

### Flags

- port: Port to launch server on. Default is `8081`.
- bucket: S3 bucket name to send files to.
- region: AWS region. Default is `us-east-1`.
