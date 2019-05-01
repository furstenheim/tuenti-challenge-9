const _ = require('lodash')
const preranges = _.times(22, function (i) {
  const r = parseInt(48 * i / 256) * 256
  return { i: i, l: 48 * i - r, u: 122 * i - r, lu: (122 * i - r) % 256, d: (122 * i - 48 * i) }
})
const ranges = _.map(preranges, function (r) {
  if (r.u !== r.lu) {
    return [[0, r.lu], [r.l, 255]]
  }
  return [[r.l, r.u]]
})
console.log(ranges)
console.log(preranges)