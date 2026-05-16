use std::fs::File;
use std::io::{self, IsTerminal, Write};
use std::path::Path;

use anyhow::{Context, Result};
use clap::ValueEnum;
use serde_json::{Map, Value};

use crate::color;

#[derive(Copy, Clone, Debug, PartialEq, Eq, ValueEnum)]
pub enum Format {
    Table,
    Json,
    Csv,
}

pub fn print(data: &[Map<String, Value>], format: Format) -> Result<()> {
    if data.is_empty() {
        println!("No results found");
        return Ok(());
    }

    match format {
        Format::Json => {
            let pretty = serde_json::to_string_pretty(data)?;
            if io::stdout().is_terminal() {
                color::print_json(&pretty);
            } else {
                println!("{pretty}");
            }
        }
        Format::Csv => {
            write_csv(io::stdout().lock(), data)?;
        }
        Format::Table => {
            print_table(data);
        }
    }
    Ok(())
}

pub fn write_to_file(
    data: &[Map<String, Value>],
    path: &str,
    format: Format,
    quiet: bool,
) -> Result<()> {
    let p = Path::new(path);
    match format {
        Format::Json => {
            let pretty = serde_json::to_string_pretty(data)?;
            std::fs::write(p, pretty).with_context(|| format!("writing {path}"))?;
        }
        Format::Csv => {
            let file = File::create(p).with_context(|| format!("creating {path}"))?;
            write_csv(file, data)?;
        }
        Format::Table => {
            anyhow::bail!("table format cannot be written to file; use --fmt json or --fmt csv");
        }
    }
    if !quiet {
        eprintln!("✓ Results saved to {path}");
    }
    Ok(())
}

fn print_table(data: &[Map<String, Value>]) {
    let keys = sorted_keys(data, 6);
    println!("\nResults ({} rows)\n", data.len());
    // header
    for k in &keys {
        print!("{:<15} ", truncate(k, 14));
    }
    println!();
    for _ in &keys {
        print!("{} ", "-".repeat(14));
    }
    println!();
    // rows
    for item in data {
        for k in &keys {
            let val = format_value(item.get(k));
            print!("{:<15} ", truncate(&val, 14));
        }
        println!();
    }
    println!();
}

pub fn print_record(record: &Map<String, Value>, title: &str) {
    println!("\n=== {title} Details ===\n");
    let mut keys: Vec<&String> = record.keys().collect();
    keys.sort();
    for k in keys {
        println!("{:<25}: {}", k, format_value(record.get(k)));
    }
    println!();
}

pub fn print_collections() {
    let collections: &[(&str, &str)] = &[
        ("artin", "Artin representations"),
        ("belyi", "Belyi maps"),
        ("char_dirichlet", "Dirichlet characters"),
        ("ec_classdata", "Elliptic curve isogeny classes"),
        ("ec_curvedata", "Elliptic curves"),
        ("g2c_curves", "Genus 2 curves"),
        ("lf_fields", "Local fields"),
        ("maass_newforms", "Maass forms"),
        ("mf_newforms", "Modular forms"),
        ("nf_fields", "Number fields"),
    ];
    println!("\nAvailable API Collections:\n");
    for (name, desc) in collections {
        println!("  {name:<20} {desc}");
    }
    println!();
}

fn sorted_keys(data: &[Map<String, Value>], max: usize) -> Vec<String> {
    let mut keys: Vec<String> = data[0].keys().cloned().collect();
    keys.sort();
    keys.truncate(max);
    keys
}

fn write_csv<W: Write>(writer: W, data: &[Map<String, Value>]) -> Result<()> {
    if data.is_empty() {
        return Ok(());
    }
    let mut wtr = csv::Writer::from_writer(writer);
    let keys: Vec<String> = {
        let mut k: Vec<String> = data[0].keys().cloned().collect();
        k.sort();
        k
    };
    wtr.write_record(&keys)?;
    for item in data {
        let row: Vec<String> = keys.iter().map(|k| format_value(item.get(k))).collect();
        wtr.write_record(&row)?;
    }
    wtr.flush()?;
    Ok(())
}

pub fn format_value(v: Option<&Value>) -> String {
    match v {
        None | Some(Value::Null) => "N/A".to_string(),
        Some(Value::String(s)) => s.clone(),
        Some(Value::Number(n)) => n.to_string(),
        Some(Value::Bool(b)) => b.to_string(),
        Some(other) => other.to_string(),
    }
}

/// UTF-8 safe truncation: cuts on character boundaries, not bytes.
pub fn truncate(s: &str, max_chars: usize) -> String {
    let count = s.chars().count();
    if count <= max_chars {
        return s.to_string();
    }
    let keep = max_chars.saturating_sub(2);
    let mut out: String = s.chars().take(keep).collect();
    out.push_str("..");
    out
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    #[test]
    fn truncate_keeps_short_strings() {
        assert_eq!(truncate("abc", 10), "abc");
    }

    #[test]
    fn truncate_handles_long_strings() {
        assert_eq!(truncate("abcdefghij", 6), "abcd..");
    }

    #[test]
    fn truncate_is_utf8_safe() {
        // 6 multibyte characters, truncate to 4 chars total => 2 chars + ".."
        let s = "αβγδεζ";
        let out = truncate(s, 4);
        assert_eq!(out.chars().count(), 4);
        assert!(out.ends_with(".."));
        assert!(out.starts_with("αβ"));
    }

    #[test]
    fn format_value_null_is_na() {
        assert_eq!(format_value(None), "N/A");
        assert_eq!(format_value(Some(&Value::Null)), "N/A");
    }

    #[test]
    fn format_value_renders_primitives() {
        assert_eq!(format_value(Some(&json!("hello"))), "hello");
        assert_eq!(format_value(Some(&json!(42))), "42");
        assert_eq!(format_value(Some(&json!(true))), "true");
    }

    #[test]
    fn sorted_keys_is_deterministic() {
        let mut m = Map::new();
        m.insert("z".into(), json!(1));
        m.insert("a".into(), json!(2));
        m.insert("m".into(), json!(3));
        let data = vec![m];
        let keys = sorted_keys(&data, 10);
        assert_eq!(keys, vec!["a", "m", "z"]);
    }
}
