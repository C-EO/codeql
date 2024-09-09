use crate::trap::TrapId;
use anyhow::Context;
use ra_ap_hir::db::DefDatabase;
use ra_ap_hir::{Crate};
use ra_ap_load_cargo::{load_workspace_at, LoadCargoConfig, ProcMacroServerChoice};
use ra_ap_project_model::CargoConfig;
use ra_ap_project_model::RustLibSource;
use std::path::{PathBuf};

mod archive;
mod config;
pub mod generated;
pub mod path;
mod translate;
pub mod trap;

fn main() -> anyhow::Result<()> {
    let cfg = config::Config::extract().context("failed to load configuration")?;
    stderrlog::new()
        .module(module_path!())
        .verbosity(2 + cfg.verbose as usize)
        .init()?;
    log::info!("{cfg:?}");
    let traps = trap::TrapFileProvider::new(&cfg).context("failed to set up trap files")?;
    let archiver = archive::Archiver {
        root: cfg.source_archive_dir,
    };

    let config = CargoConfig {
        sysroot: Some(RustLibSource::Discover),
        ..Default::default()
    };
    let progress = |t| (println!("progress: {}", t));
    let load_config = LoadCargoConfig {
        load_out_dirs_from_check: true,
        with_proc_macro_server: ProcMacroServerChoice::Sysroot,
        prefill_caches: false,
    };
    for input in cfg.inputs {
        let (db, vfs, _macro_server) = load_workspace_at(&input, &config, &load_config, &progress)
            .context("loading inputs")?;
        _macro_server.expect("no macro server");
        let crates = <dyn DefDatabase>::crate_graph(&db);
        for crate_id in crates.iter() {
            let krate = Crate::from(crate_id);
            let name = krate.display_name(&db);
            let crate_name = name
                .as_ref()
                .map(|n| n.canonical_name().as_str())
                .unwrap_or("");
            let trap = traps.create(
                "crates",
                &PathBuf::from(format!(
                    "/{}_{}",
                    crate_name,
                    crate_id.into_raw().into_u32()
                )),
            );
            translate::CrateTranslator::new(&db, trap, &krate, &vfs, &archiver)
                .emit_crate()
                .context("writing trap file")?;
        }
    }
    Ok(())
}
