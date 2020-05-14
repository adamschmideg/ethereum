const { create } = require('xmlbuilder2')

const obj = {
    root: {
        '@att': 'val',
        foo: {
            bar: 'foobar'
        },
        baz: {}
    }
};

const doc = create(obj);
const xml = doc.end({ prettyPrint: true });
console.log(xml);
