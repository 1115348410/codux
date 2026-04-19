#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

mod terminal;
mod project;

use log::{error, info};
use std::sync::Mutex;
use tauri::{Manager, State};

pub struct AppState {
    pub terminals: Mutex<Vec<terminal::Terminal>>,
    pub projects: Mutex<Vec<project::Project>>,
}

#[tauri::command]
fn create_terminal(
    working_dir: String,
    state: State<AppState>,
) -> Result<String, String> {
    info!("Creating terminal in: {}", working_dir);
    
    let mut terminals = state.terminals.lock().map_err(|e| e.to_string())?;
    
    let mut term = terminal::Terminal::new(&working_dir)
        .map_err(|e| format!("Failed to create terminal: {}", e))?;
    
    term.start();
    
    let id = term.id();
    terminals.push(term);
    
    Ok(id)
}

#[tauri::command]
fn write_terminal(
    id: String,
    input: String,
    state: State<AppState>,
) -> Result<(), String> {
    let terminals = state.terminals.lock().map_err(|e| e.to_string())?;
    
    for term in terminals.iter() {
        if term.id() == id {
            term.write(&input).map_err(|e| e.to_string())?;
            return Ok(());
        }
    }
    
    Err("Terminal not found".to_string())
}

#[tauri::command]
fn close_terminal(id: String, state: State<AppState>) -> Result<(), String> {
    let mut terminals = state.terminals.lock().map_err(|e| e.to_string())?;
    terminals.retain(|t| t.id() != id);
    Ok(())
}

#[tauri::command]
fn read_terminal_output(id: String, state: State<AppState>) -> Result<String, String> {
    let terminals = state.terminals.lock().map_err(|e| e.to_string())?;
    
    for term in terminals.iter() {
        if term.id() == id {
            return term.read_output().map_err(|e| e.to_string());
        }
    }
    
    Err("Terminal not found".to_string())
}

#[tauri::command]
fn add_project(path: String, state: State<AppState>) -> Result<project::Project, String> {
    info!("Adding project: {}", path);
    
    let project = project::Project::new(&path)?;
    
    let mut projects = state.projects.lock().map_err(|e| e.to_string())?;
    projects.push(project.clone());
    
    Ok(project)
}

#[tauri::command]
fn get_projects(state: State<AppState>) -> Result<Vec<project::Project>, String> {
    let projects = state.projects.lock().map_err(|e| e.to_string())?;
    Ok(projects.clone())
}

#[tauri::command]
fn remove_project(id: String, state: State<AppState>) -> Result<(), String> {
    let mut projects = state.projects.lock().map_err(|e| e.to_string())?;
    projects.retain(|p| p.id != id);
    Ok(())
}

fn main() {
    env_logger::init();
    
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .manage(AppState {
            terminals: Mutex::new(Vec::new()),
            projects: Mutex::new(Vec::new()),
        })
        .invoke_handler(tauri::generate_handler![
            create_terminal,
            write_terminal,
            close_terminal,
            read_terminal_output,
            add_project,
            get_projects,
            remove_project,
        ])
        .setup(|app| {
            info!("Codux starting...");
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
