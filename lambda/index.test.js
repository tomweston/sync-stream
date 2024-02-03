const { createItemFromRecord } = require('./index');

describe('createItemFromRecord', () => {
  it('should create a correct item from an S3 record', () => {
    const record = {
      s3: {
        object: {
          key: 'testFileName'
        }
      },
      eventTime: '2023-01-01T00:00:00.000Z'
    };
    const expectedItem = {
      Key: 'testFileName',
      Timestamp: '2023-01-01T00:00:00.000Z'
    };

    const item = createItemFromRecord(record);

    expect(item).toEqual(expectedItem);
  });
});
