#include <IRremote.h>

const int IRReceiverPin = 14;  // Define the IR receiver pin
IRrecv irReceiver(IRReceiverPin);

decode_results results;

void setup() {
  Serial.begin(115200);
  irReceiver.enableIRIn();  // Start the IR receiver
}

void loop() {
  if (irReceiver.decode(&results)) {
    // Check if the received IR signal is using the NEC protocol
    if (results.decode_type == NEC) {
      char receivedChar = char(results.value & 0xFF);
      Serial.print("Received Character: ");
      Serial.println(receivedChar);
    }

    irReceiver.resume();  // Receive the next IR signal
  }
}
