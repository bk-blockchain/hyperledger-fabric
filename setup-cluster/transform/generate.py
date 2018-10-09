#$PYTHON transform/generate.py --nfs-server $1 --tls-enabled $2 -o $OVERRIDE --version $VERSION --env $ENV --file $CONFIG_FILE --share $SHARE_FOLDER
#python transform/generate.py --version x86_64-1.1.0 --dbver x86_64-0.4.6 --tls-enabled false -o false --file ../driving-files/crypto-config.yaml
from string import Template
# from pathlib import Path
import string
import config as tc
import os
import sys
import argparse
import yaml


BASEDIR = os.path.dirname(__file__)
ORDERER = os.path.join(BASEDIR, "../../driving-files/crypto-config/ordererOrganizations")
PEER = os.path.join(BASEDIR, "../../driving-files/crypto-config/peerOrganizations")
KAFKA = os.path.join(BASEDIR, "../render")

#generateNamespacePod generate the yaml file to create the namespace for k8s, and return a set of paths which indicate the location of org files  

def generateKafka(DIR, override):
    tc.configKafkaNamespace(DIR, override)
    tc.configZookeepers(DIR, override)
    tc.configKafkas(DIR, override)

def generateNamespacePod(DIR, override):
	orderer0 = sorted(os.listdir(ORDERER))[0]

	orgs = []
	# remain ordered list
	for index, org in enumerate(sorted(os.listdir(DIR))):
		orgDIR = os.path.join(DIR, org) # saved dir
		## generate namespace first.
		tc.configORGS(org, orgDIR, orderer0, override, index)
		orgs.append(orgDIR)
		#orgs.append(orgDIR + "/" + DIR.lower())
	
	#print(orgs)	
	return orgs


def generateDeploymentPod(orgs, override):
	for orgindex, org in enumerate(orgs):

		if org.find("peer") != -1: #whether it create orderer pod or peer pod 
			suffix = "/peers"
		else:
			suffix = "/orderers"

		members = os.listdir(org + suffix)
		for member in members:
			#print(member)
			memberDIR = os.path.join(org + suffix, member) #saved dir
			# memberDIR = "render/"
			#print(memberDIR)
			#print(os.listdir(memberDIR))
			tc.generateYaml(member,memberDIR, suffix, override, orgindex)


#TODO kafa nodes and zookeeper nodes don't have dir to store their certificate, must use anotherway to create pod yaml.

def allInOne(override, file):
	peerOrgs = generateNamespacePod(PEER, override)
	generateDeploymentPod(peerOrgs, override)

	# check more than 1 order then run this
	stream = open(file, "r")
	YAML = yaml.load(stream)
	# if len(YAML["OrdererOrgs"]) > 1:
	if len(YAML["OrdererOrgs"][0]["Specs"]) > 0:
		generateKafka(KAFKA, override)

	ordererOrgs = generateNamespacePod(ORDERER, override)
	generateDeploymentPod(ordererOrgs, override)

def processArguments():
	parser = argparse.ArgumentParser(description='Generate network artifacts.')
	parser.add_argument('--version', dest='VERSION', type=str,
	                    help='Fabric version (default: ' + tc.VERSION + ')')
	parser.add_argument('--dbver', dest='DB_VERSION', type=str,
	                    help='Fabric version (default: ' + tc.DB_VERSION + ')')
	parser.add_argument('--kafkaver', dest='KAFKA_VERSION', type=str,
	                    help='Fabric version (default: ' + tc.KAFKA_VERSION + ')')
	parser.add_argument('--zoover', dest='ZOO_VERSION', type=str,
	                    help='Fabric version (default: ' + tc.ZOO_VERSION + ')')
	parser.add_argument('--tls-enabled', dest='TLS_ENABLED', type=str,
	                    help='Enable tls mode (default: ' + tc.TLS_ENABLED + ')')
	parser.add_argument('--env', dest='ENV', type=str,
	                    help='Fabric environment (default: ' + tc.ENV + ')')
	parser.add_argument('--file', dest='FILE', type=str,
	                    help='Config file')
	parser.add_argument("-o", "--override", dest='OVERRIDE', type=str, default="false", help="Override existing k8s yaml files")	

	# config_file = sys.argv[1] if len(sys.argv) > 1 else "cluster-config.yaml"

	# stream = open(config_file, "r")
	# YAML = yaml.load(stream)

	args = parser.parse_args()	

	tc.VERSION = args.VERSION or tc.VERSION
	tc.DB_VERSION = args.DB_VERSION or tc.DB_VERSION
	tc.KAFKA_VERSION = args.KAFKA_VERSION or tc.KAFKA_VERSION
	tc.ZOO_VERSION = args.ZOO_VERSION or tc.ZOO_VERSION
	tc.TLS_ENABLED = args.TLS_ENABLED or tc.TLS_ENABLED
	tc.ENV = args.ENV or tc.ENV

	return args

if __name__ == "__main__" :	
	args = processArguments()
	allInOne(True if args.OVERRIDE == "true" else False, args.FILE)	
	
	
	
