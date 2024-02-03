const { DynamoDBDocument } = require('@aws-sdk/lib-dynamodb');
const { DynamoDB } = require('@aws-sdk/client-dynamodb');

const dynamoDB = DynamoDBDocument.from(new DynamoDB());

const TABLE_NAME = process.env.TABLE_NAME;

function createItemFromRecord(record) {
    const s3ObjectKey = record.s3.object.key;
    const eventTime = new Date(record.eventTime).toISOString();
    return {
        Key: s3ObjectKey,
        Timestamp: eventTime
    };
}

async function insertItemIntoDynamoDB(item) {
    const params = {
        TableName: TABLE_NAME,
        Item: item
    };

    try {
        await dynamoDB.put(params);
        console.log(`Successfully inserted ${item.Key} at ${item.Timestamp}`);
    } catch (error) {
        console.error(`Error inserting into DB: ${error}`);
        throw error;
    }
}

exports.handler = async (event) => {
    for (const record of event.Records) {
        const item = createItemFromRecord(record);

        try {
            await insertItemIntoDynamoDB(item);
        } catch (error) {
            return {
                statusCode: 500,
                body: JSON.stringify('Error processing object: ' + error.message),
            };
        }
    }

    return {
        statusCode: 200,
        body: JSON.stringify('Processed: '+ s3ObjectKey),
    };
};

module.exports.createItemFromRecord = createItemFromRecord;
module.exports.insertItemIntoDynamoDB = insertItemIntoDynamoDB;