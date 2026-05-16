//! Hand-rolled ANSI syntax highlighting for JSON output.
//! Mirrors the Go implementation's color scheme.

const RESET: &str = "\x1b[0m";
const CYAN: &str = "\x1b[36m"; // keys
const GREEN: &str = "\x1b[32m"; // string values
const YELLOW: &str = "\x1b[33m"; // numbers
const MAGENTA: &str = "\x1b[35m"; // true/false/null
const GRAY: &str = "\x1b[90m"; // braces/brackets

pub fn print_json(s: &str) {
    let bytes = s.as_bytes();
    let mut out = String::with_capacity(s.len() + 64);
    let mut i = 0;

    while i < bytes.len() {
        let ch = bytes[i];
        match ch {
            b'"' => {
                // Find the closing quote (handle escapes).
                let mut j = i + 1;
                while j < bytes.len() && bytes[j] != b'"' {
                    if bytes[j] == b'\\' && j + 1 < bytes.len() {
                        j += 2;
                    } else {
                        j += 1;
                    }
                }
                let end = (j + 1).min(bytes.len());
                let token = &s[i..end];

                // Peek past whitespace to see if a `:` follows -> key.
                let mut k = end;
                while k < bytes.len() && bytes[k] == b' ' {
                    k += 1;
                }
                let is_key = k < bytes.len() && bytes[k] == b':';
                let color = if is_key { CYAN } else { GREEN };

                out.push_str(color);
                out.push_str(token);
                out.push_str(RESET);
                i = end;
            }
            b'-' | b'0'..=b'9' => {
                let mut j = i;
                while j < bytes.len() && is_number_byte(bytes[j]) {
                    j += 1;
                }
                out.push_str(YELLOW);
                out.push_str(&s[i..j]);
                out.push_str(RESET);
                i = j;
            }
            b't' if s[i..].starts_with("true") => {
                out.push_str(MAGENTA);
                out.push_str("true");
                out.push_str(RESET);
                i += 4;
            }
            b'f' if s[i..].starts_with("false") => {
                out.push_str(MAGENTA);
                out.push_str("false");
                out.push_str(RESET);
                i += 5;
            }
            b'n' if s[i..].starts_with("null") => {
                out.push_str(MAGENTA);
                out.push_str("null");
                out.push_str(RESET);
                i += 4;
            }
            b'{' | b'}' | b'[' | b']' => {
                out.push_str(GRAY);
                out.push(ch as char);
                out.push_str(RESET);
                i += 1;
            }
            _ => {
                // Preserve any multi-byte UTF-8 runs intact.
                let ch_start = i;
                let ch_len = utf8_char_len(bytes[i]);
                let ch_end = (ch_start + ch_len).min(bytes.len());
                out.push_str(&s[ch_start..ch_end]);
                i = ch_end;
            }
        }
    }
    println!("{out}");
}

fn is_number_byte(b: u8) -> bool {
    matches!(b, b'-' | b'+' | b'.' | b'e' | b'E' | b'0'..=b'9')
}

fn utf8_char_len(first: u8) -> usize {
    // ASCII or unexpected continuation byte: advance one to make progress.
    if first < 0xC0 {
        1
    } else if first < 0xE0 {
        2
    } else if first < 0xF0 {
        3
    } else {
        4
    }
}
