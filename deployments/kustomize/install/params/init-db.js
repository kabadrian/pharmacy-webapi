const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.AMBULANCE_API_MONGODB_DATABASE
const collection = process.env.AMBULANCE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
       print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })

//insert sample data
let result = db[collection].insertMany([
    {
        "id": "bobulova",
        "name": "Dr. House ambulance",
        "prescriptionList": [
            {
                "id": "predpis-01",
                "patientName": "Fero Mrkva",
                "doctorName": "Dr. House",
                "issuedDate": "2024-12-24T10:35:00Z",
                "validUntil": "2024-12-24T10:35:00Z",
                "medicines": [
                    {
                        "name": "Aspirin"
                    },
                    {
                        "name": "Lisinopril"
                    },
                    {
                        "name": "Ibuprofen"
                    }
                ]
            },
            {
                "id": "predpis-02",
                "patientName": "Jozef Golonka",
                "doctorName": "Dr. House",
                "issuedDate": "2024-12-24T10:35:00Z",
                "validUntil": "2024-12-24T10:35:00Z",
                "medicines": [
                    {
                        "name": "Aspirin"
                    },
                    {
                        "name": "Ibuprofen"
                    }
                ]
            },
            {
                "id": "predpis-01",
                "patientName": "Silvester Stalone",
                "doctorName": "Dr. House",
                "issuedDate": "2024-12-24T10:35:00Z",
                "validUntil": "2024-12-24T10:35:00Z",
                "medicines": [
                    {
                        "name": "Ibuprofen"
                    }
                ]
            }
        ],
        "medicineOrderList": [
            {
                "orderId": "order-01",
                "orderDate": "2024-05-19T12:51:58.777Z",
                "orderedBy": "Dr. House",
                "notes": "Delivery date is 3 days",
                "state": "pending",
                "medicines": [
                    {
                        "name": "Aspirin"
                    }
                ]
            },
            {
                "orderId": "order-02",
                "orderDate": "2024-05-19T12:52:23.744Z",
                "orderedBy": "Dr. House",
                "notes": "poznamka",
                "state": "pending",
                "medicines": [
                    {
                        "name": "Ibuprofen"
                    },
                    {
                        "name": "Aspirin"
                    }
                ]
            }
        ]
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
}

// exit with success
process.exit(0);