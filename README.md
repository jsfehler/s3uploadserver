# s3uploadserver
Web server that uploads files to Amazon S3

## Usage

Expects a filename, which is used for the upload.

If a form called "metadata" is included, that data will be uploaded as a separate file called <filename>_metadata.log, located next to the uploaded file.

eg: `images/myfile.png` and `images/myfile_metadata.log`

### Flags

- port: Port to launch server on. Default is 8081
- bucket: S3 bucket name to send files to
- region: AWS region. Default is us-east-1
