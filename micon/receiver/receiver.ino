#include <HTTPClient.h>
#include <IRremote.hpp>
#include <WiFi.h>
#include <WiFiMulti.h>

#define IR_RECEIVE_PIN 15

const char* ssid = "CogniSafe";
const char* password = "12345678";
// need adjustment (this is only local dev)
const char* endpoint = "http://192.168.137.1:4444/plate/verify";

const String payload = "v_id=VID8BCFBE19486F9947";

WiFiMulti wifiMulti;
void setup() {
  Serial.begin(115200);
  pinMode(2, OUTPUT);

  // Connect to Wi-Fi
  WiFi.mode(WIFI_STA);
  wifiMulti.addAP(ssid, password);
  while (wifiMulti.run() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Connecting to WiFi...");
  }
  Serial.print("Connected to WiFi in ");
  Serial.print(WiFi.localIP());
  IrReceiver.begin(IR_RECEIVE_PIN, ENABLE_LED_FEEDBACK); // Start the receiver
}

void loop() {
  if (wifiMulti.run() == WL_CONNECTED) {
    WiFiClient client;
    HTTPClient http;

    if (http.begin(client, endpoint)) {
      http.addHeader("Content-Type", "application/x-www-form-urlencoded");
      int httpResponseCode = http.POST(payload);

      if (httpResponseCode == 200) { // HTTP OK
        digitalWrite(2, HIGH); // Turn on the LED
        delay(1000); // Keep the LED on for 1 second
        digitalWrite(2, LOW); // Turn off the LED
      } else {
        Serial.print("HTTP Response code: ");
        Serial.println(httpResponseCode);
        Serial.printf("[HTTP] POST... failed, error: %s\n", http.errorToString(httpResponseCode).c_str());
      }
      
      http.end();
    } else {
       Serial.println("Client failed!");
    }
    
  }

  if (IrReceiver.decode()) {
      Serial.println(IrReceiver.decodedIRData.decodedRawData, HEX); // Print "old" raw data
      // USE NEW 3.x FUNCTIONS
      IrReceiver.printIRResultShort(&Serial); // Print complete received data in one line
      IrReceiver.printIRSendUsage(&Serial);   // Print the statement required to send this data
      IrReceiver.resume(); // Enable receiving of the next value
  }

  delay(5000); // 
}