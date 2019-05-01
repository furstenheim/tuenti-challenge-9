const ranges = [ [ [ 0, 0 ] ],
  [ [ 96, 122 ] ],
  [ [ 192, 244 ] ],
  [ [ 32, 110 ] ],
  [ [ 128, 232 ] ],
  [ [ 0, 98 ], [ 224, 255 ] ],
  [ [ 64, 220 ] ],
  [ [ 0, 86 ], [ 160, 255 ] ],
  [ [ 0, 208 ] ],
  [ [ 0, 74 ], [ 96, 255 ] ],
  [ [ 0, 255 ] ] ]

const fs = require('fs')
const input = {buf: fs.readFileSync('/dev/stdin')}

parseInput(input)
function parseInput (input) {
  const nMessages = parseInt(readLineString(input))
  for (let i = 0; i < nMessages; i++) {
    const c = parseCase(input)
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
