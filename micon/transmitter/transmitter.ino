#include <IRremote.h>

const int IRTransmitterPin = 4;  // Define the IR transmitter pin
IRsend irSender;

const char *message = "Hello";  // Message to be transmitted

void setup() {
  Serial.begin(115200);
}

void loop() {
  for (int i = 0; i < strlen(message); i++) {
    int charValue = message[i];
    irSender.sendNEC(charValue, 8);
    delay(500); 
  }

  delay(5000); 
}
