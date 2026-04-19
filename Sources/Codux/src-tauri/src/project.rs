use serde::{Deserialize, Serialize};
use std::path::Path;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Project {
    pub id: String,
    pub name: String,
    pub path: String,
}

impl Project {
    pub fn new(path: &str) -> Result<Self, String> {
        let path_obj = Path::new(path);
        
        if !path_obj.exists() {
            return Err(format!("Path does not exist: {}", path));
        }
        
        if !path_obj.is_dir() {
            return Err(format!("Path is not a directory: {}", path));
        }
        
        let name = path_obj
            .file_name()
            .and_then(|n| n.to_str())
            .unwrap_or("Unknown")
            .to_string();
        
        Ok(Self {
            id: uuid_v4(),
            name,
            path: path.to_string(),
        })
    }
}

fn uuid_v4() -> String {
    use std::time::{SystemTime, UNIX_EPOCH};
    let timestamp = SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .unwrap()
        .as_nanos();
    format!("{:x}-{:x}", timestamp, rand_u64())
}

fn rand_u64() -> u64 {
    use std::collections::hash_map::RandomState;
    use std::hash::{BuildHasher, Hasher};
    RandomState::new().build_hasher().finish()
}
