# AKFAK using Golang
### System architecture
#### Overview
---
![](screens/systemoverview.png)
---
In our project which we model after Apache Kafka, we have 4 main roles in our system.
1. Zookeeper
2. Producer
3. Broker
4. Consumer

The producer will produce messages to the Kafka system which comprises of multiple brokers and a zookeeper. The consumers in a consumer group will pull messages for a topic. 

#### Broker Setup
---
![](screens/partitionreplica.png)
---

A topic has to be created before a producer can start producing messages for that topic. During the creation of a topic, the number of partitions and replicas for those partitions will have to be specified. If a topic with 3 partitions is created, then there will be partition 0, partition 1, partition 2. And if the replication factor is 3, then for every partition there will be 3 replicas. An example of how this will be represented in the brokers is shown in the figure above.


#### Producer-Broker
---
![](screens/producerbroker.png)
---

After the topic is created, the producer can fetch the partition information for the topic and start producing messages. The producer will round-robin through the partitions. That means the first message batch it sends will be to the leader of partition-0, then next message-batch to partition-1, and so on. In this case, the partition leaders happen to be in broker-0, broker-1 and broker-2 respectively. However, since partition leaders are simply assigned to the broker that first responds to the controller's assignment, this is not always the case. When there is a distribution of leader partitions across the brokers, load-balancing amongst the brokers in handling of messages occurs.


#### Broker-Zookeeper
---
![](screens/zookeeper1.png)
---

The brokers will store the messages from the producer. Out of the brokers, there will be one broker who is the controller. The controller is in-charge of listening for the heartbeat from other brokers and pushing updates of the metadata to the other brokers. If a broker is down, it will reassign the partitions. The zookeeper will listen out for the heartbeat of the controller, and also manage the persistence of metadata information to the disk. If the controller is down, the zookeeper will elect another broker to be the new controller.


#### Broker-Consumer 
---
![](screens/brokerconsumer.png)
---

Lastly, we have consumers in consumer group. In Kafka, a consumer group can be set to subscribe to a topic, and it will assign the consumers in its group different partitions to pull from for that topic. In our implementation, we will do this manually. The user will directly specify the topic and the partition the consumer will pull from. The point of the consumer group is so that there is distribution of the pulling of messages, so that more than one machine will be responsible for pulling and processing the messages. 

### How to run
We use Docker to demonstrate how the system will work with multiple machines.
1. Download the repository using `git clone https://github.com/TAYTS/AKFAK.git`, extract it.
2. `docker-compose build` 
3. After step 2 is complete, open 3 terminals (for the minimal example of 1 producer, 1 consumer. Repeat the instructions for terminal 2 and 3 to simulate having more producers and consumers.)

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

#### Adding brokers
1. Add a new `broker_config_X.json` file in the `config` directory. Example: `config/broker_config_5.json`. Add the following into the json file, replacing `X` with the broker ID number.
```
{
  "ID": X,
  "host": "broker-X",
  "port": 5000,
  "logDir": "./logs",
  "zookeeperConnection": "zookeeper:2181"
}
```
2. Add the new broker's ID, host and port into the `data/cluster_state.json` file according to the existing conventions. Replace `X` with the broker's actual ID.
```
{
  "ID": X,
  "host": "broker-X",
  "port": 5000
}
```
3. Add the new broker information in the `docker-compose.yml` file.
    
    i) First, under `services`, add another broker. Replace `X` with the broker's actual ID.
    ```
      broker-X:
        image: akfak:latest
        container_name: broker-X
        depends_on:
          - zookeeper
        volumes:
          - brokerX-volume:/go/src/AKFAK
        command: broker -server-config ./config/broker_config_X.json
    ```
   ii) Next, under the volumes section, add the new broker's volume. Replax `X` with the broker's actual ID.
    ```
    volumes:
      brokerX-volume:
    ```

#### Testing fault-tolerance
To test for fault-tolerance, observe that the producers/consumers will be able to run smoothly as long as one broker is still alive.
| Command            | Explanation             |
| --------------     | ----------------        |
| `docker-compose rm -sf broker-x` | Simulate a broker fault |
| `docker-compose up -d broker-x ` | Revive a broker that is down |

---
#### Glossary
|Terminology            | Description|
|---                    |---|
| Producer              | Client that pushes data to the Broker |
| Consumer              | Client that pulls data from the Broker by topic and partition index|
| Broker                | Server that receives messages from producers and stores them on disk by offset. Consumers can fetch messages from the broker.
| Zookeeper             | Zookeeper serves the Broker by sending heartbeat. It is responsible for maintaining and distributing the cluster information to the broker.
| Topic                 | Topic is an identifier for producer or consumer to push or pull data pertaining to it.
| Partition             | Topics are divided into different partitions. Each partitions allow parallelization of a topic by splitting the data in a particular topic across multiple brokers
| Controller            | Unique broker that is responsible for updating metadata, reassigning partitions in the event of broker failure, etc.
| Docker                | Open source software development platform that package applications in containers, allowing them to be portable to any operating system.          |
----
