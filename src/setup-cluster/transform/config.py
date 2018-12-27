# from string import Template
from jinja2 import Environment as Template
#from pathlib import Path
import re
import string
import os

TestDir = './dest/'
PORT_START_FROM = 30500
ZOOKEEPER_PORT_START_FROM = 32750
KAFKA_PORT_START_FROM = 32730
GAP = 100  #interval for worker's port
VERSION = 'x86_64-1.1.0'
DB_VERSION = "x86_64-0.4.6"
KAFKA_VERSION = "x86_64-0.4.6"
ZOO_VERSION = "x86_64-0.4.6"
TLS_ENABLED = 'false'
ENV = "PROD"
SHARE_FOLDER = "/data"
BASEDIR = os.path.dirname(__file__)
ORDERER = os.path.join(BASEDIR, "../crypto-config/ordererOrganizations")
PEER = os.path.join(BASEDIR, "../crypto-config/peerOrganizations")
KAFKA = os.path.join(BASEDIR, "../crypto-config/kafka")
NODE_1 = "ip-172-31-44-159"
NODE_2 = "ip-172-31-39-152"
NODE_3 = "ip-172-31-38-202"

def render(src, dest, **kw):
        # t = Template(open(src, 'r').read())
        t = Template(
    line_statement_prefix='%',
    variable_start_string="${",
    variable_end_string="}"
        ).from_string(open(src, 'r').read())
        options = dict(
                version = VERSION,
                tlsEnabled = TLS_ENABLED,
                dbVersion = DB_VERSION,
                kafkaVersion = KAFKA_VERSION,
                zooVersion = ZOO_VERSION,
                **kw)
        with open(dest, 'w') as f:
                f.write(t.render(**options))
                # f.write(t.substitute(**options))

        ##### For testing ########################
        ##testDest = dest.split("/")[-1]        ##
        ##with open(TestDir+testDest, 'w') as d:##
        ##d.write(t.substitute(**kw))           ##
        ##########################################

def condRender(src, dest, override, **kw):
  if not os.path.exists(dest):
      render(src, dest, **kw)
  elif os.path.exists(dest) and override:
      render(src, dest, **kw)

def getTemplate(templateName):
        baseDir = os.path.dirname(__file__)
        configTemplate = os.path.join(baseDir, "../templates/" + templateName)
        return configTemplate

def getAddressSegment(index):
        # pattern = re.compile('(\d+)$')
        # result = re.search(pattern, name.split("-")[0])
        # return (int(result.group(0)) -1 if result else 0) * GAP
        return index * GAP

def configKafkaNamespace(path, override):
    namespaceTemplate = getTemplate("template-kafka-namespace.yaml")
    condRender(namespaceTemplate, path + "/" + "kafka-namespace.yaml", override)

# bydefault 3 kafka and 4 zookeeper as channel, and multiple orderer will be scale based on this
def configZookeepers(path, override):
    for i in range(0, 3):
        zkTemplate = getTemplate("template-zookeeper.yaml")
        zkPodName = "zookeeper" + str(i) + "-kafka"
        zookeeperID = "zookeeper" + str(i)
        seq = i + 1
        nodePort1 = ZOOKEEPER_PORT_START_FROM + (i * 3 + 1)
        nodePort2 = nodePort1 + 1
        nodePort3 = nodePort2 + 1
        zooServersTemplate = "server.1=zookeeper0.kafka:2888:3888 server.2=zookeeper1.kafka:2888:3888 server.3=zookeeper2.kafka:2888:3888"
        zooServers = zooServersTemplate.replace("zookeeper" + str(i) + ".kafka", "0.0.0.0")
        hostname=NODE_1

        condRender(zkTemplate, path + "/" + zookeeperID + "-kafka.yaml", override,
           zkPodName=zkPodName,
           zookeeperID=zookeeperID,
           seq=seq,
           zooServers=zooServers,
           nodePort1=nodePort1,
           nodePort2=nodePort2,
           nodePort3=nodePort3,
           hostname=hostname
                                )


def configKafkas(path, override):
    for i in range(0, 4):
        kafkaTemplate = getTemplate("template-kafka.yaml")
        kafkaPodName = "kafka" + str(i) + "-kafka"
        kafkaID = "kafka" + str(i)
        seq = i
        nodePort1 = KAFKA_PORT_START_FROM + (i * 2 + 1)
        nodePort2 = nodePort1 + 1
        advertisedHostname = "kafka" + str(i) + ".kafka"
        hostname = NODE_1

        condRender(kafkaTemplate, path + "/" + kafkaID + "-kafka.yaml", override,
           kafkaPodName=kafkaPodName,
           kafkaID=kafkaID,
           seq=seq,
           advertisedHostname=advertisedHostname,
           nodePort1=nodePort1,
           nodePort2=nodePort2,
           hostname=hostname
        )



# create org/namespace
# copy to SHARE_FOLDER => need to map to nfs
def configORGS(name, path, orderer0, override, index): # name means if of org, path describe where is the namespace yaml to be created.

        hostPath = path.replace("transform/../../", SHARE_FOLDER + "/")
        # hostPath = path.replace("transform/../../", "/data/")

        if path.find("peer") != -1 :
                ####### pod config yaml for org cli
                cliTemplate = getTemplate("template-cli.yaml")

                mspPathTemplate = 'users/Admin@{}/msp'
                tlsPathTemplate =  'users/Admin@{}/tls'

                hostname = ''
                if ("1" in name):
                    hostname = NODE_2
                elif ('2' in name):
                    hostname = NODE_3
                # path + "/" + name
                condRender(cliTemplate, "render" + "/" + name + "-cli.yaml", override,
                        podName = "cli",
                        namespace = name,
                        mspPath = mspPathTemplate.format(name),
                        tlsPath = tlsPathTemplate.format(name),
                        corePeerID = "peer0." + name,
                        peerAddress = "peer0." + name + ":7051",
                        mspid = name.split('-')[0].capitalize()+"MSP",
                        hostname=hostname,
                )
                #######

                ####### pod config yaml for org ca

                ###Need to expose pod's port to worker ! ####
                ##org format like this org1-f-1##
                # addressSegment = (int(name.split("-")[0].split("org")[-1]) - 1) * GAP
                addressSegment = getAddressSegment(index)
                # each oganization should have unique ip, so ip + port should be unique
                exposedPort = PORT_START_FROM + addressSegment

                caTemplate = getTemplate("template-ca.yaml")

                tlsCertTemplate = '/etc/hyperledger/fabric-ca-server-config/{}-cert.pem'
                tlsKeyTemplate = '/etc/hyperledger/fabric-ca-server-config/{}'
                caPathTemplate = 'ca/'

                cmdTemplate =  ' fabric-ca-server start --ca.certfile /etc/hyperledger/fabric-ca-server-config/{}-cert.pem --ca.keyfile /etc/hyperledger/fabric-ca-server-config/{} -b admin:adminpw -d '

                skFile = ""
                for f in os.listdir(path+"/ca"):  # find out sk!
                        if f.endswith("_sk"):
                                skFile = f

                condRender(caTemplate, "render" + "/" + name + "-ca.yaml", override,
                        namespace = name,
                        command = '"' + cmdTemplate.format("ca."+name, skFile) + '"',
                        caPath = caPathTemplate,
                        tlsKey = tlsKeyTemplate.format(skFile),
                        tlsCert = tlsCertTemplate.format("ca."+name),
                        nodePort = exposedPort,
                        path = hostPath,
                        hostname=hostname,
                )
                #######

def generateYaml(member, memberPath, flag, override, index):
        if flag == "/peers":
                configPEERS(member, memberPath, override, index)
        else:
                configORDERERS(member, memberPath, override, index)


# create peer/pod
def configPEERS(name, path, override, index): # name means peerid.
        configTemplate = getTemplate("template-peer.yaml")
        hostPath = path.replace("transform/../../", SHARE_FOLDER + "/")
        # hostPath = path.replace("transform/../", "render/")
        mspPathTemplate = 'peers/{}/msp'
        tlsPathTemplate =  'peers/{}/tls'
        #mspPathTemplate = './msp'
        #tlsPathTemplate = './tls'
        nameSplit = name.split(".")
        peerName = nameSplit[0]
        orgName = nameSplit[1]

        # addressSegment = (int(orgName.split("-")[0].split("org")[-1]) - 1) * GAP
        addressSegment = getAddressSegment(index)
        ##peer from like this peer 0##
        peerOffset = int((peerName.split("peer")[-1])) * 4
        exposedPort1 = PORT_START_FROM + addressSegment + peerOffset + 1
        exposedPort2 = PORT_START_FROM + addressSegment + peerOffset + 2
        exposedPort3 = PORT_START_FROM + addressSegment + peerOffset + 3
        exposedPort4 = PORT_START_FROM + addressSegment + peerOffset + 4

        hostname = ''
        if ("1" in orgName):
            hostname = NODE_2
        elif ('2' in orgName):
            hostname = NODE_3

        # path + "/" + name
        condRender(configTemplate, "render/" + name + ".yaml", override,
                namespace = orgName,
                podName = peerName + "-" + orgName,
                peerID  = peerName,
                org = orgName,
                corePeerID = name,
    # peerAddress and peerCCAddress are for chaincode container to connect
                peerAddress = name + ":7051",
        # peerAddress = name + ":" + str(exposedPort1),
                peerCCAddress = "0.0.0.0" + ":7052",
        # peerCCAddress = name + ":" + str(exposedPort2),
                peerGossip = name  + ":7051",
                localMSPID = orgName.split('-')[0].capitalize()+"MSP",
                mspPath = mspPathTemplate.format(name),
                tlsPath = tlsPathTemplate.format(name),
                nodePort1 = exposedPort1,
                nodePort2 = exposedPort2,
                nodePort3 = exposedPort3,
                nodePort4 = exposedPort4,
                path=hostPath,
                hostname=hostname,
        # version 1.0, 0.6 will not using address auto detect
        addressAutoDetect = "false" if re.match(r"^(?:1\.0|0\.6)\.*", VERSION) else "true",
        peerCmd = "start --peer-chaincodedev=true" if ENV == "DEV" else "start"
        )


# create orderer/pod
def configORDERERS(name, path, override, index): # name means ordererid
        configTemplate = getTemplate("template-orderer.yaml")
        hostPath = path.replace("transform/../../", SHARE_FOLDER + "/")
        # hostPath = path.replace("transform/../", "render/")
        # genesisPath = os.path.dirname(os.path.dirname(hostPath))
        genesisPath = os.path.dirname(SHARE_FOLDER + "/driving-files/channel-artifacts/")
        mspPathTemplate = 'orderers/{}/msp'
        tlsPathTemplate = 'orderers/{}/tls'

        nameSplit = name.split(".")
        ordererName = nameSplit[0]
        orgName = nameSplit[1]
        ordererOffset = int(ordererName.split("orderer")[-1])
        addressSegment = getAddressSegment(index)
        exposedPort = 32000 + addressSegment + ordererOffset

        hostname = ''
        if ("orderer0" in ordererName):
            hostname = NODE_2
        elif ("orderer1" in ordererName):
            hostname = NODE_3
        #path + "/" + name
        condRender(configTemplate, "render/" + name + ".yaml", override,
                namespace = orgName.lower(),
                ordererID = ordererName,
        #       podName =  ordererName + "-" + orgName,
                podName =  ordererName,
                localMSPID =  orgName.capitalize() + "MSP",
                nodePort = exposedPort,
                path = hostPath,
                genesis = genesisPath,
                hostname=hostname,
        )


#if __name__ == "__main__":
#       #ORG_NUMBER = 3
#       podFile = Path('./fabric_cluster.yaml')
#       if podFile.is_file():
#               os.remove('./fabric_cluster.yaml')

#delete the previous exited file
#       configPeerORGS(1, 2)
#       configPeerORGS(2, 2)
#       configOrdererORGS()
