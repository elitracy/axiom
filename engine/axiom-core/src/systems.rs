#[derive(Clone, Copy)]
pub enum Status {
    Online,
    Degraded,
    Critical,
}

pub trait Subsystem {
    fn tick(&mut self);
    fn get_id(&self) -> &str;
    fn get_name(&self) -> &str;
    fn get_health(&self) -> f64;
    fn get_status(&self) -> Status;
    fn get_degradation_rate(&self) -> f64;
}

pub struct PowerSystem {
    id: String,
    name: String,
    health: f64,
    status: Status,
    degradation_rate: f64,
}

impl Subsystem for PowerSystem {
    fn tick(&mut self) {
        self.health -= 0.05;
    }

    fn get_id(&self) -> &str {
        &self.id
    }

    fn get_name(&self) -> &str {
        &self.name
    }

    fn get_health(&self) -> f64 {
        self.health
    }

    fn get_status(&self) -> Status {
        self.status
    }

    fn get_degradation_rate(&self) -> f64 {
        self.degradation_rate
    }
}
