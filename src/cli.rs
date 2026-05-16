use anyhow::{Context, Result};
use clap::{Args, Parser, Subcommand};

use crate::output::{self, Format};
use crate::query::{self, EcOptions, NfOptions};

pub const VERSION: &str = env!("CARGO_PKG_VERSION");

#[derive(Parser)]
#[command(
    name = "lmfdb",
    version = VERSION,
    about = "Query the LMFDB (L-Functions and Modular Forms Database) from the command line",
    disable_version_flag = false,
)]
pub struct Cli {
    #[command(subcommand)]
    pub command: Command,
}

#[derive(Subcommand)]
pub enum Command {
    /// Query number fields
    Nf(NfArgs),
    /// Query elliptic curves
    Ec(EcArgs),
    /// List available API collections
    #[command(alias = "ls")]
    List,
    /// Show version information
    #[command(alias = "v")]
    Version,
    /// Install Chrome browser for reCAPTCHA bypass
    InstallBrowser,
}

#[derive(Args)]
pub struct NfArgs {
    /// Number field degree
    #[arg(short = 'd', long, default_value_t = 2)]
    pub degree: u32,
    /// Filter by discriminant
    #[arg(long)]
    pub disc: Option<String>,
    /// Filter by class number
    #[arg(long = "class")]
    pub class_num: Option<String>,
    /// Filter by signature (e.g. "0,1")
    #[arg(long = "sig")]
    pub signature: Option<String>,
    /// Number of results
    #[arg(short = 'n', long = "limit", default_value_t = 10)]
    pub limit: u32,
    /// Result offset for pagination
    #[arg(long, default_value_t = 0)]
    pub offset: u32,
    /// Sort by field (prefix `-` for descending)
    #[arg(long)]
    pub sort: Option<String>,
    /// Fields to return (comma-separated)
    #[arg(short = 'f', long = "fields")]
    pub fields: Option<String>,
    /// Output file
    #[arg(short = 'o', long = "output")]
    pub output: Option<String>,
    /// Output format
    #[arg(long = "fmt", value_enum, default_value_t = Format::Table)]
    pub format: Format,
    /// Quiet mode
    #[arg(short = 'q', long)]
    pub quiet: bool,
    /// Get specific field by label (e.g. 2.0.3.1)
    #[arg(long)]
    pub id: Option<String>,
    /// Use headless browser (bypasses reCAPTCHA)
    #[arg(long)]
    pub browser: bool,
}

#[derive(Args)]
pub struct EcArgs {
    /// Filter by Mordell-Weil rank
    #[arg(short = 'r', long)]
    pub rank: Option<String>,
    /// Filter by torsion order
    #[arg(short = 't', long)]
    pub torsion: Option<String>,
    /// Filter by conductor
    #[arg(long)]
    pub conductor: Option<String>,
    /// Number of results
    #[arg(short = 'n', long = "limit", default_value_t = 10)]
    pub limit: u32,
    /// Result offset for pagination
    #[arg(long, default_value_t = 0)]
    pub offset: u32,
    /// Sort by field (prefix `-` for descending)
    #[arg(long)]
    pub sort: Option<String>,
    /// Fields to return (comma-separated)
    #[arg(short = 'f', long = "fields")]
    pub fields: Option<String>,
    /// Output file
    #[arg(short = 'o', long = "output")]
    pub output: Option<String>,
    /// Output format
    #[arg(long = "fmt", value_enum, default_value_t = Format::Table)]
    pub format: Format,
    /// Quiet mode
    #[arg(short = 'q', long)]
    pub quiet: bool,
    /// Use headless browser (bypasses reCAPTCHA)
    #[arg(long)]
    pub browser: bool,
}

pub fn run(cli: Cli) -> Result<()> {
    match cli.command {
        Command::Nf(args) => run_nf(args),
        Command::Ec(args) => run_ec(args),
        Command::List => {
            output::print_collections();
            Ok(())
        }
        Command::Version => {
            println!("LMFDB CLI v{VERSION}");
            Ok(())
        }
        Command::InstallBrowser => {
            anyhow::bail!("`install-browser` is not implemented in the Rust port yet (chromiumoxide will be added under a feature flag)")
        }
    }
}

fn run_nf(args: NfArgs) -> Result<()> {
    let opts = NfOptions {
        degree: args.degree,
        disc: args.disc,
        class_num: args.class_num,
        signature: args.signature,
        limit: args.limit,
        offset: args.offset,
        sort: args.sort,
        fields: args.fields,
        id: args.id,
        browser: args.browser,
    };
    let url = query::build_nf_url(&opts);
    if !args.quiet {
        eprintln!("Querying LMFDB...");
        eprintln!("  {url}");
    }
    let data = query::fetch(&url, opts.browser).context("fetching from LMFDB")?;

    if opts.id.is_some() && data.len() == 1 {
        if let Some(path) = args.output.as_deref() {
            output::write_to_file(&data, path, args.format, args.quiet)?;
        } else {
            output::print_record(&data[0], "Number Field");
        }
        return Ok(());
    }

    if let Some(path) = args.output.as_deref() {
        output::write_to_file(&data, path, args.format, args.quiet)?;
    } else {
        output::print(&data, args.format)?;
    }
    Ok(())
}

fn run_ec(args: EcArgs) -> Result<()> {
    let opts = EcOptions {
        rank: args.rank,
        torsion: args.torsion,
        conductor: args.conductor,
        limit: args.limit,
        offset: args.offset,
        sort: args.sort,
        fields: args.fields,
        browser: args.browser,
    };
    let url = query::build_ec_url(&opts);
    if !args.quiet {
        eprintln!("Querying LMFDB...");
        eprintln!("  {url}");
    }
    let data = query::fetch(&url, opts.browser).context("fetching from LMFDB")?;

    if let Some(path) = args.output.as_deref() {
        output::write_to_file(&data, path, args.format, args.quiet)?;
    } else {
        output::print(&data, args.format)?;
    }
    Ok(())
}
