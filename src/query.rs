use std::time::Duration;

use anyhow::{anyhow, bail, Result};
use serde_json::{Map, Value};
use url::Url;

const BASE_URL: &str = "https://www.lmfdb.org";
const USER_AGENT: &str = concat!(
    "lmfdb-cli/",
    env!("CARGO_PKG_VERSION"),
    " (+https://github.com/FrankieeW/lmfdb-cli)"
);

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
    let base = Url::parse(BASE_URL).expect("BASE_URL is valid");

    if let Some(id) = &opt.id {
        let mut url = base.join("/api/nf_fields/").expect("static path");
        url.path_segments_mut()
            .expect("URL has path segments")
            .pop_if_empty()
            .push(id)
            .push("");
        url.query_pairs_mut().append_pair("_format", "json");
        return url.into();
    }

    let mut url = base.join("/api/nf_fields/").expect("static path");
    {
        let mut q = url.query_pairs_mut();
        q.append_pair("_format", "json");
        q.append_pair("_limit", &opt.limit.to_string());
        q.append_pair("degree", &format!("i{}", opt.degree));
        if opt.offset > 0 {
            q.append_pair("_offset", &opt.offset.to_string());
        }
        if let Some(sort) = &opt.sort {
            q.append_pair("_sort", sort);
        }
        if let Some(disc) = &opt.disc {
            q.append_pair("disc", &format!("i{disc}"));
        }
        if let Some(cn) = &opt.class_num {
            q.append_pair("class_number", &format!("i{cn}"));
        }
        if let Some(sig) = &opt.signature {
            q.append_pair("signature", &format!("li{}", sig.replace(',', ";")));
        }
        if let Some(fields) = &opt.fields {
            q.append_pair("_fields", fields);
        }
    }
    url.into()
}

pub fn build_ec_url(opt: &EcOptions) -> String {
    let mut url = Url::parse(BASE_URL)
        .expect("BASE_URL is valid")
        .join("/api/ec_curvedata/")
        .expect("static path");
    {
        let mut q = url.query_pairs_mut();
        q.append_pair("_format", "json");
        q.append_pair("_limit", &opt.limit.to_string());
        if opt.offset > 0 {
            q.append_pair("_offset", &opt.offset.to_string());
        }
        if let Some(sort) = &opt.sort {
            q.append_pair("_sort", sort);
        }
        if let Some(r) = &opt.rank {
            q.append_pair("rank", &format!("i{r}"));
        }
        if let Some(t) = &opt.torsion {
            q.append_pair("torsion", &format!("i{t}"));
        }
        if let Some(c) = &opt.conductor {
            q.append_pair("conductor", c);
        }
        if let Some(fields) = &opt.fields {
            q.append_pair("_fields", fields);
        }
    }
    url.into()
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

    if !status.is_success() {
        bail!(
            "HTTP {status}: {}",
            body.chars().take(200).collect::<String>()
        );
    }

    let lower = body.to_ascii_lowercase();
    if lower.contains("recaptcha") || lower.contains("checking your browser") {
        bail!("blocked by reCAPTCHA (try --browser; not yet ported)");
    }

    parse_records(&body)
}

fn parse_records(body: &str) -> Result<Vec<Map<String, Value>>> {
    let value: Value =
        serde_json::from_str(body).map_err(|e| anyhow!("parsing JSON response: {e}"))?;

    // Two response shapes:
    //   { "data": [ {...}, ... ] }  -- list endpoint
    //   { ... }                      -- single record endpoint
    if let Some(data) = value.get("data") {
        if let Some(arr) = data.as_array() {
            let mut out = Vec::with_capacity(arr.len());
            let mut skipped = 0usize;
            for item in arr {
                if let Some(obj) = item.as_object() {
                    out.push(obj.clone());
                } else {
                    skipped += 1;
                }
            }
            if skipped > 0 {
                eprintln!("warning: skipped {skipped} non-object record(s) in response");
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
    fn nf_id_url_includes_format() {
        let mut o = nf_default();
        o.id = Some("2.0.3.1".into());
        o.degree = 5; // ignored when id is set
        let url = build_nf_url(&o);
        assert_eq!(
            url,
            "https://www.lmfdb.org/api/nf_fields/2.0.3.1/?_format=json"
        );
    }

    #[test]
    fn nf_id_url_encodes_path_segment() {
        let mut o = nf_default();
        // A pathological id with a slash and a query separator.
        o.id = Some("foo/bar?baz".into());
        let url = build_nf_url(&o);
        // The slash and ? in the id must be percent-encoded, not interpreted as
        // path/query separators.
        assert!(url.starts_with("https://www.lmfdb.org/api/nf_fields/"));
        assert!(url.contains("foo%2Fbar%3Fbaz"));
        assert!(url.ends_with("?_format=json"));
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
        assert!(url.starts_with("https://www.lmfdb.org/api/nf_fields/?"));
        assert!(url.contains("_limit=10"));
        assert!(url.contains("degree=i2"));
        assert!(url.contains("_offset=20"));
        assert!(url.contains("_sort=-class_number"));
        assert!(url.contains("disc=i-5"));
        assert!(url.contains("class_number=i1"));
        assert!(url.contains("signature=li0%3B1"));
        // commas are unreserved per RFC 3986 and may pass through.
        assert!(url.contains("_fields=label") && url.contains("disc"));
    }

    #[test]
    fn nf_url_encodes_special_chars_in_filters() {
        let mut o = nf_default();
        // `&` and `=` would corrupt the query if not encoded.
        o.disc = Some("1&x=2".into());
        let url = build_nf_url(&o);
        assert!(url.contains("disc=i1%26x%3D2"));
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
        assert!(url.contains("rank=i2"));
        assert!(url.contains("torsion=i5"));
        assert!(url.contains("conductor=11"));
        assert!(url.contains("_sort=conductor"));
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
