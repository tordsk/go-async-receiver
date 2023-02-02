package main

import (
	"context"
	"encoding/json"
	"fmt"
	pr "github.com/crow-misia/go-push-receiver"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

type EventLogger struct {
	logrus.FieldLogger
}

func (receiver EventLogger) Output(calldepth int, s string) error {
	receiver.Infof("%d: %s", calldepth, s)
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	done := make(chan bool, 1)
	go func() {
		<-c
		signal.Stop(c)
		cancel()
		done <- true
	}()
	l := logrus.New()
	l.SetLevel(logrus.TraceLevel)
	logger := logrus.NewEntry(l)

	opts := []pr.ClientOption{
		pr.WithHeartbeatPeriod(10 * time.Second),
		pr.WithLogger(EventLogger{logger}),
	}

	creds, err := loadCreds()
	if err != nil {
		logger.WithError(err).Error()
	} else {
		opts = append(opts, pr.WithCreds(creds))
	}

	client := pr.New("730815416529", opts...)
	go client.Subscribe(ctx)

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				logger.WithError(err).Error()
			}
			return
		case event := <-client.Events:
			logger.WithField("event", event).
				WithField("type", fmt.Sprintf("%T\n", event)).
				Infof("received notification event")
			switch e := event.(type) {
			case *pr.UpdateCredentialsEvent:
				err := persist(e.Credentials)
				logger.Info(e.Credentials.Token)
				if err != nil {
					logger.WithError(err).Error()
				}
			case *pr.ConnectedEvent:
				logger.WithField("connected", e.ServerTimestamp).Info("Connected")
			case *pr.DisconnectedEvent:
				logger.WithError(e.ErrorObj).Error("disconnected")
			case *pr.HeartbeatEvent:
				logger.WithField("heartbeat", e).Debug()
			case *pr.MessageEvent:
				logger.WithField("msg", e).Info(string(e.Data))
			case *pr.RetryEvent:
				logger.WithField("retry_after", e.RetryAfter).WithError(e.ErrorObj).Error("retry")
			}
		}
	}
}

func persist(credentials *pr.FCMCredentials) error {
	if credentials == nil {
		return fmt.Errorf("credentials is nil")
	}
	content, err := json.MarshalIndent(*credentials, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile("credentials.json", content, 0666)
}

func loadCreds() (*pr.FCMCredentials, error) {
	content, err := os.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}
	creds := pr.FCMCredentials{}
	return &creds, json.Unmarshal(content, &creds)
}

/**
Tesla: 729118971144 | KIA: 1014617823644
[x] Need to find the senderID for the vendor
sender id was extracted from the decompiled app. Should do similar for other apps awell. Unique per vendor


[x] Need a way to register our firebase token (app) to the user
Should be a request to register device token

ENABLE: Gotta be one of:
First:
GET: https://owner-api.teslamotors.com/api/1/notification_preferences?device_type=android-production-firebase&locale=en_GB&app_version=4.17.0-1566&platform=android&device_token=fXMw9HoVTs-esibt5MDedI%3AAPA91bFud9nZywYJkKLoMwqcLFn18nIECgBgAxc6XyFnrMq3q5fQ0aazgassbQId_h_ExyhTdtzo_ufxXLiFBi0ev7ssPaanShFAbwYQFnMLPPIT_CebCbRObWSnz6X7vgbMwSvTK83f

2.
POST: https://owner-api.teslamotors.com/api/1/users/keys
{
  "public_key": "04a6ae11068713d4ef96ea5101a9c633aad867aa2b7d11262b00659d4010a508f56dfa94964a679a1b9f4c9717752ff5e2ef8b105aea3595f065a5d77dbf679000",
  "name": "RePhone 1",
  "kind": "mobile_device",
  "battery_level": "unknown",
  "bluetooth_enabled": true,
  "device_token": "fXMw9HoVTs-esibt5MDedI:APA91bFud9nZywYJkKLoMwqcLFn18nIECgBgAxc6XyFnrMq3q5fQ0aazgassbQId_h_ExyhTdtzo_ufxXLiFBi0ev7ssPaanShFAbwYQFnMLPPIT_CebCbRObWSnz6X7vgbMwSvTK83f",
  "low_power_mode": "unknown",
  "locale": "en_GB",
  "model": "Nothing A063",
  "tag": "TeslaApp/4.17.0-1566/043f9a8590/android/12",
  "bluetooth_authorization": "authorized",
  "device_type": "android-production-firebase",
  "location_permissions": "notDetermined"
}
3.
POST: https://owner-api.teslamotors.com/api/1/notification_preferences
{
  "device_token": "fXMw9HoVTs-esibt5MDedI:APA91bFud9nZywYJkKLoMwqcLFn18nIECgBgAxc6XyFnrMq3q5fQ0aazgassbQId_h_ExyhTdtzo_ufxXLiFBi0ev7ssPaanShFAbwYQFnMLPPIT_CebCbRObWSnz6X7vgbMwSvTK83f",
  "locale": "en_GB",
  "device_type": "android-production-firebase",
  "app_version": "4.17.0-1566",
  "platform": "android",
  "notification_preferences": {
    "charging_started": false,
    "charging_interrupted": true,
    "bml_complete_order": true,
    "sentry_off_no_ap": true,
    "vpp_event_beginning_discharge": true,
    "rewrap_vault_keys": true,
    "cabin_overheat_protection": false,
    "update_available": true,
    "wait_for_user_no_inverters_ready": true,
    "vehicle_driving_in_valet_mode": true,
    "climate_keeper_ended_soc": true,
    "key_removed": true,
    "vehicle_added": true,
    "service_rideshare_credits_available": true,
    "internal_app_update_nag": true,
    "suspicious_activity": true,
    "inbox": true,
    "vpp_event_scheduled": true,
    "service_installables_appointment": true,
    "storm_mode_on": true,
    "energy_support_chat_agent_joined_chat": true,
    "removed_authorized_client": true,
    "vehicle_twelve_volt_battery_replacement_alert": true,
    "urgent_can_alert": true,
    "wait_for_jump_start": true,
    "update_imminent": true,
    "autopark_forward_started": false,
    "energy_support_chat_new_message": true,
    "grid_outage": true,
    "service_survey": true,
    "service_loaner_vehicle_to_be_provided": true,
    "questionnaire": true,
    "enrollment_notification_rejected": true,
    "wait_for_user_low_soe": true,
    "service_complete": true,
    "custom_energy_alert": true,
    "notifications_preconditioning_complete": true,
    "service_tracker_reminder": true,
    "autopark_completed_success": false,
    "expired_payment": true,
    "alarm": true,
    "service_in_service_with_estimated_completion_time": true,
    "enrollment_notification_ineligible": true,
    "wait_for_user_retries_exhausted": true,
    "key_added": true,
    "service_range_analysis_complete": true,
    "service_rideshare_credits_to_be_provided": true,
    "service_loaner_vehicle_accept_agreement": true,
    "service_estimate_available_reminder": true,
    "lootbox_credit_bonus": true,
    "wait_for_solar": true,
    "autopark_unavailable_plugged_in": false,
    "service_estimate_available": true,
    "off_grid_approaching_low_soe": true,
    "service_appointment_reminder": true,
    "black_start_failure": true,
    "preconditioning_complete": true,
    "service_vehicle_self_test_request": true,
    "service_ready_for_pickup": true,
    "service_in_part_hold_with_estimated_completion_time": true,
    "service_range_analysis_incomplete": true,
    "power_rationality_alert": true,
    "lootbox_quarter_lottery": true,
    "vehicle_unsecure": true,
    "supercharging_disabled": true,
    "climate_off_timeout": true,
    "scheduled_update_failed_to_start": true,
    "tpms_alert": true,
    "outstanding_balance": true,
    "charging_complete": false,
    "please_move_car": true,
    "low_soe": true,
    "speed_limit_proximity_triggered": true,
    "service_manual_message": true,
    "scheduled_island_contactor_open": true,
    "service_in_service": true,
    "service_new_chat_message": true,
    "service_questions_prompt": true,
    "supercharging_entity_switch": true,
    "supercharging_survey": true,
    "sentry_panic": true,
    "energy_support_chat_agent_ended_chat": true,
    "lootbox_credit_expiration": true,
    "charge_pricing_information": true,
    "incentive_inspection": true,
    "factory_reset_initiated": true,
    "refer_friend": true,
    "enrollment_notification_participating": true,
    "climate_keeper_warning": true,
    "climate_keeper_critical": true,
    "dog_mode_faulted": true,
    "added_authorized_client": true,
    "climate_keeper_reminder": true,
    "sentry_off_soc": true,
    "service_in_part_hold": true,
    "lootbox_promotion": true,
    "lootbox_chargingmiles_expiration": true,
    "car_active": true,
    "sentry_on_extended": true,
    "charge_cable_unlatched": true,
    "battery_breaker_open": true,
    "climate_keeper_ended_fault": true,
    "secret_level": true,
    "high_usage_supercharger": true,
    "climate_ended": true
  }
}



DISABLE:
https://owner-api.teslamotors.com/api/1/device/{device_token}/deactivate

*/
