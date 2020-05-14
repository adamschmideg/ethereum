const { create } = require('xmlbuilder2')
const assert = require('assert')

it('xml', () => {
    const obj = {
        root: {
            arr: [
                {foo: "well"},
                {bar: 32},
                "alone",
                {"@id": 2},
                {"@id": 42, nested: 42},
                {"@id": 3},
            ],
        }
    };

    const doc = create(obj);
    const xml = doc.end({ prettyPrint: true });
    const wanted = `
<?xml version="1.0"?>
<root att="val">
  <foo>
    <bar>foobar</bar>
  </foo>
  <baz/>
</root>
    `
    assert.deepEqual(xml, '')
 //    expect(xml).toBe('')
//    expect(xml).toBe(wanted)
})

