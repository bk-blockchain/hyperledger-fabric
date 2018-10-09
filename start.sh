cd driving-files
./generateArtifacts.sh
cd ..
./copy-files.sh
cd setup-cluster
python transform/generate.py --version x86_64-1.1.0 --dbver x86_64-0.4.6 --tls-enabled true -o false --file ../driving-files/crypto-config.yaml
cd ..
./deploy.sh 

