const int LOCK = 7;
const int DELAY = 8000;
const int BAUD = 19200;

void setup() {
  Serial.begin(BAUD);
  pinMode(LOCK, OUTPUT);
}

void loop() {
  while (Serial.available() > 0) {
    int val = Serial.parseInt();
    if (val == 1) {
      unlock();
    }
  }
}

void unlock() {
  digitalWrite(LOCK, HIGH);
  delay(DELAY);
  digitalWrite(LOCK, LOW);
}
