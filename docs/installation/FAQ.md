# FAQ

## 1. Data Loss Problem

- Kappital-Manager is deployed in Kubernetes Cluster, these two components will use the Kubernetes’s Schedule function. When these two components scheduler from one node to the other or restart, **the database’s data will lose**. Because these data is saving in container, it will lose with the re-scheduler and restart.
- If the actual situation needs data persistence, please read **FAQ 3 or Deploy the kappital in physical machine**.

## 2. How to recover after restarting container

Recovering data after restarting container, please use the following steps:
1. Using `docker cp` copy the database and package files from container to local host.
2. Restart the container.
3. Using `docker cp` copy the database and package files from local host to the restarted container.

## 3. How to fix lost data problem as deploying with container

Some security hardening is performed during container deployment, for example, non-root users. If you want to use the default database and ensure data reliability, you need to deploy the database as the root user and mount the database to a physical machine node. However, if Kubernetes schedules related services, data will be lost because Kubernetes cannot migrate related data.

The service package must be run as the root user to ensure the data reliability of the service package.

P.S: If change to the other database please read **FAQ 4** to extend the database interface to connect to other databases, such as MySQL.

## 4. Expand Interface Problem

If you want to change the database mode, refer to the following file to implement related interfaces:
- Database Interface: [db.go](../../pkg/models/db.go)

## 5. How to Create Certificates

Here only offer one method to create the certificate, if this certificate is not satisfied your using situation, please change to the other certificate creating method.

```shell
openssl genrsa -out ca.key 2048
openssl req -new -out ca.csr -key ca.key -subj "/C=CN/CN=Kappital - RootCA" -keyform PEM
openssl x509 -req -in ca.csr -out ca.crt -signkey ca.key  -CAcreateserial -days 365

openssl genrsa -out client.key 2048
openssl req -new -out client.csr -key client.key -subj "/C=CN/CN=Kappital - Client" -keyform PEM
echo "subjectAltName = IP:{alt name for the ip address}" > extfile-client.cnf
echo "extendedKeyUsage=clientAuth" >> extfile-client.cnf
openssl x509 -req -in client.csr -out client.crt -CA ca.crt -CAkey ca.key -CAcreateserial -days 365 -extfile extfile-client.cnf

openssl genrsa -out server.key 2048
openssl req -new -out server.csr -key server.key -subj "/C=CN/CN=Kappital - Server" -keyform PEM
echo "subjectAltName = IP:{alt name for the ip address}" > extfile-server.cnf
echo "extendedKeyUsage=serverAuth" >> extfile-server.cnf
openssl x509 -req -in server.csr -out server.crt -CA ca.crt -CAkey ca.key  -CAcreateserial -days 365 -extfile extfile-server.cnf
```


