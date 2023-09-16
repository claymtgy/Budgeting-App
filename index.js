const readline = require('readline').createInterface({
  input: process.stdin,
  output: process.stdout
});

readline.question('Who are you? \n', name => {
  console.log(`Hey there ${name}!`);
  readline.close();
});
