/* Automatically generated nanopb constant definitions */
/* Generated by nanopb-0.3.6-dev */

#include "config_msg.pb.h"

/* @@protoc_insertion_point(includes) */
#if PB_PROTO_HEADER_VERSION != 30
#error Regenerate this file with the current version of nanopb generator.
#endif

const char Wifi_ssid_default[33] = "";
const char Wifi_password_default[64] = "";
const uint16_t APA102_frameRate_default = 60u;
const uint16_t APA102_numLights_default = 150u;
const uint32_t APA102_SPIspeed_default = 8000000u;
const bool DisplaySettings_enabled_default = false;
const uint32_t DisplaySettings_I2Cspeed_default = 4000000u;
const char Host_name_default[32] = "";
const bool Host_highSpeed_default = false;
const char Config_romURL_default[128] = "";


const pb_field_t Wifi_fields[3] = {
    PB_FIELD(  1, STRING  , OPTIONAL, STATIC  , FIRST, Wifi, ssid, ssid, &Wifi_ssid_default),
    PB_FIELD(  2, STRING  , OPTIONAL, STATIC  , OTHER, Wifi, password, ssid, &Wifi_password_default),
    PB_LAST_FIELD
};

const pb_field_t APA102_fields[4] = {
    PB_FIELD(  1, UINT32  , OPTIONAL, STATIC  , FIRST, APA102, frameRate, frameRate, &APA102_frameRate_default),
    PB_FIELD(  2, UINT32  , OPTIONAL, STATIC  , OTHER, APA102, numLights, frameRate, &APA102_numLights_default),
    PB_FIELD(  3, UINT32  , OPTIONAL, STATIC  , OTHER, APA102, SPIspeed, numLights, &APA102_SPIspeed_default),
    PB_LAST_FIELD
};

const pb_field_t DisplaySettings_fields[3] = {
    PB_FIELD(  1, BOOL    , OPTIONAL, STATIC  , FIRST, DisplaySettings, enabled, enabled, &DisplaySettings_enabled_default),
    PB_FIELD(  2, UINT32  , OPTIONAL, STATIC  , OTHER, DisplaySettings, I2Cspeed, enabled, &DisplaySettings_I2Cspeed_default),
    PB_LAST_FIELD
};

const pb_field_t Host_fields[3] = {
    PB_FIELD(  1, STRING  , OPTIONAL, STATIC  , FIRST, Host, name, name, &Host_name_default),
    PB_FIELD(  2, BOOL    , OPTIONAL, STATIC  , OTHER, Host, highSpeed, name, &Host_highSpeed_default),
    PB_LAST_FIELD
};

const pb_field_t Config_fields[7] = {
    PB_FIELD(  1, MESSAGE , OPTIONAL, STATIC  , FIRST, Config, wifiClient, wifiClient, &Wifi_fields),
    PB_FIELD(  2, MESSAGE , OPTIONAL, STATIC  , OTHER, Config, wifiAP, wifiClient, &Wifi_fields),
    PB_FIELD(  3, MESSAGE , OPTIONAL, STATIC  , OTHER, Config, apa102, wifiAP, &APA102_fields),
    PB_FIELD(  4, MESSAGE , OPTIONAL, STATIC  , OTHER, Config, host, apa102, &Host_fields),
    PB_FIELD(  5, MESSAGE , OPTIONAL, STATIC  , OTHER, Config, display, host, &DisplaySettings_fields),
    PB_FIELD(  6, STRING  , OPTIONAL, STATIC  , OTHER, Config, romURL, display, &Config_romURL_default),
    PB_LAST_FIELD
};


/* Check that field information fits in pb_field_t */
#if !defined(PB_FIELD_32BIT)
/* If you get an error here, it means that you need to define PB_FIELD_32BIT
 * compile-time option. You can do that in pb.h or on compiler command line.
 * 
 * The reason you need to do this is that some of your messages contain tag
 * numbers or field sizes that are larger than what can fit in 8 or 16 bit
 * field descriptors.
 */
PB_STATIC_ASSERT((pb_membersize(Config, wifiClient) < 65536 && pb_membersize(Config, wifiAP) < 65536 && pb_membersize(Config, apa102) < 65536 && pb_membersize(Config, host) < 65536 && pb_membersize(Config, display) < 65536), YOU_MUST_DEFINE_PB_FIELD_32BIT_FOR_MESSAGES_Wifi_APA102_DisplaySettings_Host_Config)
#endif

#if !defined(PB_FIELD_16BIT) && !defined(PB_FIELD_32BIT)
/* If you get an error here, it means that you need to define PB_FIELD_16BIT
 * compile-time option. You can do that in pb.h or on compiler command line.
 * 
 * The reason you need to do this is that some of your messages contain tag
 * numbers or field sizes that are larger than what can fit in the default
 * 8 bit descriptors.
 */
PB_STATIC_ASSERT((pb_membersize(Config, wifiClient) < 256 && pb_membersize(Config, wifiAP) < 256 && pb_membersize(Config, apa102) < 256 && pb_membersize(Config, host) < 256 && pb_membersize(Config, display) < 256), YOU_MUST_DEFINE_PB_FIELD_16BIT_FOR_MESSAGES_Wifi_APA102_DisplaySettings_Host_Config)
#endif


/* @@protoc_insertion_point(eof) */
