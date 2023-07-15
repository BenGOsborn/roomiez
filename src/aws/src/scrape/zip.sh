mkdir -p nodejs
cp node_modules package* nodejs -r
zip /mnt/c/Users/bengr/Downloads/layer.zip nodejs -r
zip /mnt/c/Users/bengr/Downloads/lambda.zip *.js
rm -rf nodejs