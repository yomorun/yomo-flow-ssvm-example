use std::io::{self, Read};
use std::str::FromStr;
use ssvm_wasi_helper::ssvm_wasi_helper::_initialize;

fn main() {
	_initialize();
	let mut buffer = String::new();
	io::stdin().read_to_string(&mut buffer).expect("Error reading from STDIN");
	let args: Vec<String> = serde_json::from_str(&buffer).unwrap();
	let x = f64::from_str(&args[0]).unwrap();
	print!("{:?}", 3.0 * x);
}
