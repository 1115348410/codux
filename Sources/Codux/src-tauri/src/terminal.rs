use log::{error, info};
use portable_pty::{native_pty_system, CommandBuilder, PtySize, PtyPair, Child};
use std::collections::HashMap;
use std::io::{Read, Write};
use std::sync::mpsc::{channel, Receiver, Sender};
use std::thread;
use std::sync::{Arc, Mutex};

pub struct Terminal {
    id: String,
    working_dir: String,
    pty_pair: PtyPair,
    child: Box<dyn Child + Send + Sync>,
    output_buffer: Arc<Mutex<String>>,
    #[allow(dead_code)]
    reader_handle: thread::JoinHandle<()>,
}

impl Terminal {
    pub fn new(working_dir: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let pty_system = native_pty_system();
        
        let pair = pty_system.openpty(PtySize {
            rows: 24,
            cols: 80,
            pixel_width: 0,
            pixel_height: 0,
        })?;
        
        let mut cmd = CommandBuilder::new("cmd.exe");
        cmd.cwd(working_dir);
        
        let child = pair.slave.spawn_command(cmd)?;
        
        let pty_pair = pty_system.openpty(PtySize {
            rows: 24,
            cols: 80,
            pixel_width: 0,
            pixel_height: 0,
        })?;
        
        let output_buffer = Arc::new(Mutex::new(String::new()));
        let buffer_clone = output_buffer.clone();
        
        let mut reader = pair.master.take_reader()?;
        
        let reader_handle = thread::spawn(move || {
            let mut buf = [0u8; 1024];
            loop {
                match reader.read(&mut buf) {
                    Ok(0) => break,
                    Ok(n) => {
                        let text = String::from_utf8_lossy(&buf[..n]).to_string();
                        let mut buffer = buffer_clone.lock().unwrap();
                        buffer.push_str(&text);
                    }
                    Err(e) => {
                        error!("Terminal read error: {}", e);
                        break;
                    }
                }
            }
        });
        
        Ok(Self {
            id: uuid_v4(),
            working_dir: working_dir.to_string(),
            pty_pair: pair,
            child,
            output_buffer,
            reader_handle,
        })
    }
    
    pub fn id(&self) -> String {
        self.id.clone()
    }
    
    pub fn write(&self, input: &str) -> Result<(), Box<dyn std::error::Error>> {
        let mut writer = self.pty_pair.master.take_writer()?;
        writer.write_all(input.as_bytes())?;
        writer.write_all(b"\r\n")?;
        writer.flush()?;
        Ok(())
    }
    
    pub fn read_output(&self) -> Result<String, Box<dyn std::error::Error>> {
        let buffer = self.output_buffer.lock().map_err(|e| e.to_string())?;
        Ok(buffer.clone())
    }
    
    pub fn start(&mut self) {
        info!("Terminal {} started in {}", self.id, self.working_dir);
    }
}

impl Drop for Terminal {
    fn drop(&mut self) {
        let _ = self.child.kill();
        info!("Terminal {} closed", self.id);
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
