# AKFAK

## Distributed System Course Project

### How to run
We use docker to demonstrate how the system will work with multiple machines.
1. Download the repository using `git clone https://github.com/TAYTS/AKFAK.git`, extract it.
2. `docker-compose build` 
3. After step 2 is complete, open 3 terminals.

#### On the first terminal, we will set up our zookeeper and brokers.
| Command            | Explanation             |
| --------------     | ----------------        |
|`docker-compose up` | This will set up the the zookeeper and 4 brokers. |

As there are 4 brokers, in the commands below, `X` in `<broker-X:port>` can be substituted by any number from \[1-4]. Port is specified in the `broker_X_config.json` files to be `5000`.

#### On the second terminal, we will create a topic and start a producer instance that will produce messages for a specified topic.
| Command            | Explanation             |
| --------------     | ----------------        |
|`docker container run --rm  -it --network=kafka-net akfak bash`| This runs a bash terminal in a container in the network that is set up. |
| `admin-topic -kafka-server <broker-X:port> create -topic <topic_name> -partitions <num_partitions> -replica-factor <num_replicas>` | Creates a topic with `<topic_name>` with `<num_partitions>` partitions and `<num_replicas>` replicas. <br>**Example**: `admin-topic -kafka-server broker-1:5000 create -topic topic1 -partitions 3 -replica-factor 3` |
|`producer --kafka-server broker-X:port --topic <topic_name>` | This instantiates a producer. The producer will call the given broker address to retrieve partition information on the topic it is producing.<br>**Example**: `producer --kafka-server broker-1:5000 --topic topic1`|

#### On the third terminal, we will create a consumer instance and pull data from a specified topic and partition.
| Command            | Explanation             |
| --------------     | ----------------        |
|`docker container run --rm  -it --network=kafka-net akfak bash` | This runs a bash terminal in a container in the network that is set up. |
|`consumer --id <id> --kafka-server <broker-X:port> --topic <topic_name>` | This instantiates a consumer. The consumer will call the given broker address to retrieve partition information on the topic it is producing. The user will then select one partition the consumer will pull from.<br>**Example**: `consumer --id 1 --kafka-server broker-1:5000 --topic topic1`|

