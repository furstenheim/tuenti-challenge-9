module.exports = {
  splitMessage,
  parseInput,
  valueInRanges,
  findMinimumInterstitialLength,
  findCombination
}
const _ = require('lodash')
const allRanges = [ [ [ 0, 0 ] ],
  [ [ 48, 122 ] ],
  [ [ 96, 244 ] ],
  [ [ 0, 110 ], [ 144, 255 ] ],
  [ [ 0, 255 ] ]]
const preranges = [ { i: 0, l: 0, u: 0, lu: 0, d: 0 },
  { i: 1, l: 48, u: 122, lu: 122, d: 74 },
  { i: 2, l: 96, u: 244, lu: 244, d: 148 },
  { i: 3, l: 144, u: 366, lu: 110, d: 222 },
  { i: 4, l: 192, u: 488, lu: 232, d: 296 }]

// parseInput(input)
function parseInput (input) {
  const nMessages = parseInt(readLineString(input))
  for (let i = 0; i < nMessages; i++) {
    const c = parseCase(input)
  }
}

function solveCase (c) {

}

function findCombinationMod (value, nCharacters) {

}

function findCombination (value, nCharacters) {
  const prerange = preranges[nCharacters]
  if (!(prerange.l <= value && prerange.u >= value) ) {
    throw new Error('Impossible combination')
  }
  const combination = []
  for (let i = nCharacters - 1; i >= 0; i--) {
    const char = Math.min(122, value - i * 48)
    combination[i] = char
    value -= char
  }
  return combination
}


function setUpComputations (c) {
  const original = c.original
  const copy = c.copy
  const splitOriginal = splitMessage(c.original)
  const splitCopy = splitMessage(c.copy)
  const {diff, payloadHash} = findMinimumInterstitialLength(splitOriginal, splitCopy)
}

function findMinimumInterstitialLength (splitOriginal, splitCopy) {
  const payloadHash = []
  for (let diff = 0; diff < allRanges.length * 16; diff++) {
    let isPossible = true
    for (let i = 0; i < 16; i++) {
      const repetitions = i < (diff % 16) ? parseInt(diff / 16) + 1 : (parseInt(diff / 16))
      const ranges = allRanges[repetitions]
      if (!ranges) {
        throw new Error(`Bad index: ${repetitions}`)
      }
      const hashPosition = (i + splitCopy.m1.length) % 16
      const vOriginal = splitOriginal.hash[hashPosition]
      const vCopy = splitCopy.m1Hash[hashPosition] + splitCopy.m2UnshiftedHash[((hashPosition - splitCopy.m1.length - diff) % 16 + 16) % 16]
      const necessaryValue = ((vOriginal - vCopy) % 256 + 256) % 256
      payloadHash[i] = necessaryValue
      const inRange = valueInRanges(necessaryValue, ranges)
      if (!inRange) {
        isPossible = false
        break
      }
    }
    if (isPossible) {
      return {diff, payloadHash}
    }
  }
  throw new Error('Diff not possible')


}
function valueInRanges (v, ranges) {
  return !!_.find(ranges, function (range) {
    return range[0] <= v && range[1] >= v
  })

}
function splitMessage (m) {
  const mFirstSplitIndex = m.indexOf('---')
  const m1 = m.slice(0, mFirstSplitIndex + 3)
  const m2 = m.slice(mFirstSplitIndex + 3)
  const m1Hash = _.times(16, () => 0)
  const m2Hash = _.times(16, () => 0)
  const m2UnshiftedHash = _.times(16, () => 0)
  _.forEach(m1, function (c, i) {
    m1Hash[i % 16] = (m1Hash[i % 16] + c) % 256
  })
  _.forEach(m2, function (c, j) {
    const i = j + m1.length
    m2Hash[i % 16] = (m2Hash[i % 16] + c) % 256
    m2UnshiftedHash[j % 16] = (m2UnshiftedHash[j % 16] + c) % 256
  })
  return {
    m1Hash,
    m1,
    m2Hash,
    m2,
    m2UnshiftedHash,
    hash: _.zipWith(m1Hash, m2Hash, (h1, h2) => (h1 + h2) % 256)
  }
}
function parseCase (input) {
  const nLinesOriginal = parseInt(readLineString(input))
  if (!nLinesOriginal) {
    throw new Error('Input is empty')
  }

  const original = readNLines(input, nLinesOriginal)
  if (original.includes('\n') || original.includes('\r')) {
    throw new Error(`Original contains new line: ${original}`)
  }
  const nLinesCopy = parseInt(readLineString(input))
  if (!nLinesCopy) {
    throw new Error('Copy is empty')
  }
  const copy = readNLines(input, nLinesCopy)
  if (copy.includes('\n') || copy.includes('\r')) {
    throw new Error(`Copy contains new line: ${copy}`)
  }
  return {
    original,
    copy
  }
}

function readLineString (input) {
  const buf = input.buf
  const newLineIndex = buf.indexOf('\n')
  input.buf = buf.slice(newLineIndex + 1)
  return buf.slice(0, newLineIndex).toString()
}
function readLineBuffer (input) {
  const buf = input.buf
  const newLineIndex = buf.indexOf('\n')
  input.buf = buf.slice(newLineIndex + 1)
  return buf.slice(0, newLineIndex)
}
function readNLines (input, n) {
  let b = Buffer.from('')
  for (let i = 0; i < n; i++) {
    b = Buffer.concat([b, readLineBuffer(input)])
  }
  return b
}
