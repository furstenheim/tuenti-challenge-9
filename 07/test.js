// same input and output
// less than 16 char
// one less the other one more
const _ = require('lodash')
const lib = require('./lib')
const {expect} = require('chai')
console.log('###########')
describe('Split message', function () {
  const tests = [{
    m: Buffer.from('Subject: Boat;From: Charlie;To: Desmond;------Not Penny\'s boat'),
    hash: [122, 103, 95, 92, -123 + 256, -89 + 256, -78 + 256, 14, 44, -8 + 256, 99, 56, 86, 75, -50 + 256, 1]
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const o = lib.splitMessage(t.m)
      expect(o.m1Hash).lengthOf(16)
      expect(o.m2Hash ).lengthOf(16)

      console.log(o)
      const hash = _.zipWith(o.m1Hash, o.m2Hash, (h1, h2) => (h1 + h2) % 256)
      console.log(hash)
      expect(hash).deep.equal(t.hash)
    })
  })
})
