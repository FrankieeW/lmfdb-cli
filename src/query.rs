use std::time::Duration;

use anyhow::{anyhow, bail, Result};
use serde_json::{Map, Value};

const BASE_URL: &str = "https://www.lmfdb.org";
const USER_AGENT: &str =
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 lmfdb-cli";

pub struct NfOptions {
    pub degree: u32,
    pub disc: Option<String>,
    pub class_num: Option<String>,
    pub signature: Option<String>,
    pub limit: u32,
    pub offset: u32,
    pub sort: Option<String>,
    pub fields: Option<String>,
    pub id: Option<String>,
    pub browser: bool,
}

pub struct EcOptions {
    pub rank: Option<String>,
    pub torsion: Option<String>,
    pub conductor: Option<String>,
    pub limit: u32,
    pub offset: u32,
    pub sort: Option<String>,
    pub fields: Option<String>,
    pub browser: bool,
}

pub fn build_nf_url(opt: &NfOptions) -> String {
    if let Some(id) = &opt.id {
        return format!("{BASE_URL}/api/nf_fields/{id}/?_format=json");
    }

    let mut url = format!(
        "{BASE_URL}/api/nf_fields/?_format=json&_limit={}&degree=i{}",
        opt.limit, opt.degree
    );

    if opt.offset > 0 {
        url.push_str(&format!("&_offset={}", opt.offset));
    }
    if let Some(sort) = &opt.sort {
        url.push_str("&_sort=");
        url.push_str(sort);
    }
    if let Some(disc) = &opt.disc {
        url.push_str("&disc=i");
        url.push_str(disc);
    }
    if let Some(cn) = &opt.class_num {
        url.push_str("&class_number=i");
        url.push_str(cn);
    }
    if let Some(sig) = &opt.signature {
        url.push_str("&signature=li");
        url.push_str(&sig.replace(',', ";"));
    }
    if let Some(fields) = &opt.fields {
        url.push_str("&_fields=");
        url.push_str(fields);
    }
    url
}

pub fn build_ec_url(opt: &EcOptions) -> String {
    let mut url = format!(
        "{BASE_URL}/api/ec_curvedata/?_format=json&_limit={}",
        opt.limit
    );

    if opt.offset > 0 {
        url.push_str(&format!("&_offset={}", opt.offset));
    }
    if let Some(sort) = &opt.sort {
        url.push_str("&_sort=");
        url.push_str(sort);
    }
    if let Some(r) = &opt.rank {
        url.push_str("&rank=i");
        url.push_str(r);
    }
    if let Some(t) = &opt.torsion {
        url.push_str("&torsion=i");
        url.push_str(t);
    }
    if let Some(c) = &opt.conductor {
        url.push_str("&conductor=");
        url.push_str(c);
    }
    if let Some(fields) = &opt.fields {
        url.push_str("&_fields=");
        url.push_str(fields);
    }
    url
}

pub fn fetch(url: &str, browser: bool) -> Result<Vec<Map<String, Value>>> {
    if browser {
        bail!("--browser is not yet implemented in the Rust port; use the API path or the Go binary as a fallback");
    }
    fetch_api(url)
}

fn fetch_api(url: &str) -> Result<Vec<Map<String, Value>>> {
    let client = reqwest::blocking::Client::builder()
        .timeout(Duration::from_secs(60))
        .user_agent(USER_AGENT)
        .gzip(true)
        .build()?;

    let resp = client.get(url).send()?;
    let status = resp.status();
    let body = resp.text()?;

    if body.contains("recaptcha") || body.contains("Checking your browser") {
        bail!("blocked by reCAPTCHA (try --browser; not yet ported)");
    }
    if !status.is_success() {
        bail!("HTTP {status}: {}", body.chars().take(200).collect::<String>());
    }

    parse_records(&body)
}

fn parse_records(body: &str) -> Result<Vec<Map<String, Value>>> {
    let value: Value = serde_json::from_str(body)
        .map_err(|e| anyhow!("parsing JSON response: {e}"))?;

    // Two response shapes:
    //   { "data": [ {...}, ... ] }  -- list endpoint
    //   { ... }                      -- single record endpoint
    if let Some(data) = value.get("data") {
        if let Some(arr) = data.as_array() {
            let mut out = Vec::with_capacity(arr.len());
            for item in arr {
                if let Some(obj) = item.as_object() {
                    out.push(obj.clone());
                }
            }
            return Ok(out);
        }
    }
    if let Some(obj) = value.as_object() {
        return Ok(vec![obj.clone()]);
    }
    Ok(Vec::new())
}

#[cfg(test)]
mod tests {
    use super::*;

    fn nf_default() -> NfOptions {
        NfOptions {
            degree: 2,
            disc: None,
            class_num: None,
            signature: None,
            limit: 10,
            offset: 0,
            sort: None,
            fields: None,
            id: None,
            browser: false,
        }
    }

    #[test]
    fn nf_id_url_skips_filters() {
        let mut o = nf_default();
        o.id = Some("2.0.3.1".into());
        o.degree = 5; // ignored when id is set
        let url = build_nf_url(&o);
        assert_eq!(url, "https://www.lmfdb.org/api/nf_fields/2.0.3.1/?_format=json");
    }

    #[test]
    fn nf_list_url_composes_filters() {
        let mut o = nf_default();
        o.disc = Some("-5".into());
        o.class_num = Some("1".into());
        o.signature = Some("0,1".into());
        o.offset = 20;
        o.sort = Some("-class_number".into());
        o.fields = Some("label,disc".into());
        let url = build_nf_url(&o);
        assert!(url.starts_with("https://www.lmfdb.org/api/nf_fields/?_format=json&_limit=10&degree=i2"));
        assert!(url.contains("&_offset=20"));
        assert!(url.contains("&_sort=-class_number"));
        assert!(url.contains("&disc=i-5"));
        assert!(url.contains("&class_number=i1"));
        assert!(url.contains("&signature=li0;1"));
        assert!(url.contains("&_fields=label,disc"));
    }

    #[test]
    fn ec_url_composes_filters() {
        let opt = EcOptions {
            rank: Some("2".into()),
            torsion: Some("5".into()),
            conductor: Some("11".into()),
            limit: 25,
            offset: 0,
            sort: Some("conductor".into()),
            fields: None,
            browser: false,
        };
        let url = build_ec_url(&opt);
        assert!(url.contains("_limit=25"));
        assert!(url.contains("&rank=i2"));
        assert!(url.contains("&torsion=i5"));
        assert!(url.contains("&conductor=11"));
        assert!(url.contains("&_sort=conductor"));
    }

    #[test]
    fn parse_records_handles_list_shape() {
        let body = r#"{"data":[{"label":"2.0.3.1","disc":-3},{"label":"2.0.4.1","disc":-4}]}"#;
        let records = parse_records(body).unwrap();
        assert_eq!(records.len(), 2);
        assert_eq!(records[0]["label"], "2.0.3.1");
    }

    #[test]
    fn parse_records_handles_single_object() {
        let body = r#"{"label":"2.0.3.1","disc":-3}"#;
        let records = parse_records(body).unwrap();
        assert_eq!(records.len(), 1);
        assert_eq!(records[0]["disc"], -3);
    }
}
