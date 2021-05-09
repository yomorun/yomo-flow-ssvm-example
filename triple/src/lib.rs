use wasm_bindgen::prelude::*;

#[wasm_bindgen]
pub fn triple(x:f64) -> f64 {
	return 3.0 * x;
}

