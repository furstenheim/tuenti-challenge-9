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

describe('findMinimumInterstitialLength', function () {
  const tests = [{
    original: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Not Penny's boat`),
    copy: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Penny's boat :)`),
    expected: {diff: 49, payloadHash: [ 124, 39, 75, 33, 31, 15, 178, 226, 236, 71, 251, 89, 197, 4, 191, 238 ]}
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const originalSplit = lib.splitMessage(t.original)
      const copySplit = lib.splitMessage(t.copy)
      console.log(originalSplit)
      console.log(copySplit)
      const l = lib.findMinimumInterstitialLength(originalSplit, copySplit)
      expect(l).deep.equal(t.expected)
    })
  })
})

describe.only('solveCase', function () {
  const tests = [{
    original: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Not Penny's boat`),
    copy: Buffer.from(`Subject: Boat;From: Charlie;To: Desmond;------Penny's boat :)`),
    expected: '03W000000S0e0000Xzzwue08BzQz0Z0DzzzzzzRzzzzzez_zz'
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const l = lib.solveCase(t)
      expect(l).deep.equal(t.expected)
    })
  })
})

describe('Find combination', function () {
  const tests = [{
    nCharacters: 0,
    value: 0,
    expected: []
  }, {
    nCharacters: 1,
    value: 48,
    expected: [48]
  }, {
    nCharacters: 1,
    value: 49,
    expected: [49]
  }, {
    nCharacters: 2,
    value: 98,
    expected: [48, 50]
  }, {
    nCharacters: 3,
    value: 250,
    expected: [48, 80, 122]
  }, {
    nCharacters: 4,
    value: 250,
    expected: [48, 48, 48, 106]
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const combination = lib.findCombination(t.value, t.nCharacters)
      expect(combination).deep.equal(t.expected)
    })
  })
})
describe('Find combination mod', function () {
  const tests = [{
    nCharacters: 0,
    value: 0,
    expected: []
  }, {
    nCharacters: 1,
    value: 48,
    expected: [48]
  }, {
    nCharacters: 1,
    value: 49,
    expected: [49]
  }, {
    nCharacters: 2,
    value: 98,
    expected: [48, 50]
  }, {
    nCharacters: 3,
    value: 250,
    expected: [48, 80, 122]
  }, {
    nCharacters: 4,
    value: 250,
    expected: [48, 48, 48, 106]
  }, {
    nCharacters: 4,
    value: 257,
    expected: [48, 48, 48, 113]
  }, {
    nCharacters: 4,
    value: 193,
    expected: [48, 48, 48, 49]
  }]
  _.forEach(tests, function (t, i) {
    it('#' + i, function () {
      const combination = lib.findCombination(t.value, t.nCharacters)
      expect(combination).deep.equal(t.expected)
    })
  })
})
