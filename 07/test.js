// same input and output
// less than 16 char
// one less the other one more
// tests minimal with same length
const _ = require('lodash')
const lib = require('./lib')
const {expect} = require('chai')
describe('Split message', function () {
  const tests = [{
    m: Buffer.from('Subject: Boat;From: Charlie;To: Desmond;------Not Penny\'s boat'),
    hash: [122, 103, 95, 92, -123 + 256, -89 + 256, -78 + 256, 14, 44, -8 + 256, 99, 56, 86, 75, -50 + 256, 1]
  },{
    m: Buffer.from('------'),
    hash: [45, 45, 45, 45, 45, 45, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
  }]

  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const o = lib.splitMessage(t.m)
      expect(o.m1Hash).lengthOf(16)
      expect(o.m2Hash).lengthOf(16)

      console.log(o)
      const hash = _.zipWith(o.m1Hash, o.m2Hash, (h1, h2) => (h1 + h2) % 256)
      console.log(hash)
      expect(hash).deep.equal(t.hash)
    })
  })
})
describe('Value in ranges', function () {
  const tests = [{
    r: [[0, 220]],
    v: 1,
    expected: true
  }, {
    r: [[0, 220]],
    v: -1,
    expected: false
  }, {
    r: [[0, 220]],
    v: 221,
    expected: false
  }, {
    r: [[0, 220]],
    v: 0,
    expected: true
  }, {
    r: [[0, 220], [230, 255]],
    v: 231,
    expected: true
  }, {
    r: [[0, 220], [230, 255]],
    v: 225,
    expected: false
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const inside = lib.valueInRanges(t.v, t.r)
      expect(inside).equal(t.expected)
    })
  })
})

describe.only('findMinimumInterstitialLength', function () {
  const tests = [{
    original: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Not Penny's boat`),
    copy: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Penny's boat :)`),
    expected: 49
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const originalSplit = lib.splitMessage(t.original)
      const copySplit = lib.splitMessage(t.copy)
      console.log(originalSplit)
      console.log(copySplit)
      const l = lib.findMinimumInterstitialLength(originalSplit, copySplit)
      expect(l).equal(t.expected)
    })
  })

})