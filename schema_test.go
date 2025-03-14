package parquet_test

import (
	"testing"
	"time"

	"github.com/segmentio/parquet-go"
)

func TestSchemaOf(t *testing.T) {
	tests := []struct {
		value interface{}
		print string
	}{
		{
			value: new(struct{ Name string }),
			print: `message {
	required binary Name (STRING);
}`,
		},

		{
			value: new(struct {
				X int
				Y int
			}),
			print: `message {
	required int64 X (INT(64,true));
	required int64 Y (INT(64,true));
}`,
		},

		{
			value: new(struct {
				X float32
				Y float32
			}),
			print: `message {
	required float X;
	required float Y;
}`,
		},

		{
			value: new(struct {
				Inner struct {
					FirstName string `parquet:"first_name"`
					LastName  string `parquet:"last_name"`
				} `parquet:"inner,optional"`
			}),
			print: `message {
	optional group inner {
		required binary first_name (STRING);
		required binary last_name (STRING);
	}
}`,
		},

		{
			value: new(struct {
				Short float32 `parquet:"short,split"`
				Long  float64 `parquet:"long,split"`
			}),
			print: `message {
	required float short;
	required double long;
}`,
		},

		{
			value: new(struct {
				Inner struct {
					FirstName          string `parquet:"first_name"`
					ShouldNotBePresent string `parquet:"-"`
				} `parquet:"inner,optional"`
			}),
			print: `message {
	optional group inner {
		required binary first_name (STRING);
	}
}`,
		},

		{
			value: new(struct {
				Inner struct {
					FirstName    string `parquet:"first_name"`
					MyNameIsDash string `parquet:"-,"`
				} `parquet:"inner,optional"`
			}),
			print: `message {
	optional group inner {
		required binary first_name (STRING);
		required binary - (STRING);
	}
}`,
		},

		{
			value: new(struct {
				Inner struct {
					TimestampMillis int64 `parquet:"timestamp_millis,timestamp"`
					TimestampMicros int64 `parquet:"timestamp_micros,timestamp(microsecond)"`
				} `parquet:"inner,optional"`
			}),
			print: `message {
	optional group inner {
		required int64 timestamp_millis (TIMESTAMP(isAdjustedToUTC=true,unit=MILLIS));
		required int64 timestamp_micros (TIMESTAMP(isAdjustedToUTC=true,unit=MICROS));
	}
}`,
		},

		{
			value: new(struct {
				Name string `parquet:",json"`
			}),
			print: `message {
	required binary Name (JSON);
}`,
		},

		{
			value: new(struct {
				A map[int64]string `parquet:"," parquet-key:",timestamp"`
				B map[int64]string
			}),
			print: `message {
	required group A (MAP) {
		repeated group key_value {
			required int64 key (TIMESTAMP(isAdjustedToUTC=true,unit=MILLIS));
			required binary value (STRING);
		}
	}
	required group B (MAP) {
		repeated group key_value {
			required int64 key (INT(64,true));
			required binary value (STRING);
		}
	}
}`,
		},

		{
			value: new(struct {
				A map[int64]string `parquet:",optional" parquet-value:",json"`
			}),
			print: `message {
	optional group A (MAP) {
		repeated group key_value {
			required int64 key (INT(64,true));
			required binary value (JSON);
		}
	}
}`,
		},

		{
			value: new(struct {
				A map[int64]string `parquet:",optional"`
			}),
			print: `message {
	optional group A (MAP) {
		repeated group key_value {
			required int64 key (INT(64,true));
			required binary value (STRING);
		}
	}
}`,
		},

		{
			value: new(struct {
				A map[int64]string `parquet:",optional" parquet-value:",json" parquet-key:",timestamp(microsecond)"`
			}),
			print: `message {
	optional group A (MAP) {
		repeated group key_value {
			required int64 key (TIMESTAMP(isAdjustedToUTC=true,unit=MICROS));
			required binary value (JSON);
		}
	}
}`,
		},

		{
			value: new(struct {
				Ushort uint16 `parquet:"ushort"`
				Short  int16  `parquet:"short"`
			}),
			print: `message {
	required int32 ushort (INT(16,false));
	required int32 short (INT(16,true));
}`,
		},

		{
			value: new(struct {
				IntDate *int32 `parquet:"int_date,date"`
			}),
			print: `message {
	optional int32 int_date (DATE);
}`,
		},

		{
			value: new(struct {
				IntTime    *int64     `parquet:"int_time,timestamp"`
				TimeTime   *time.Time `parquet:"time_time,timestamp"`
				NotUTCTime *time.Time `parquet:"not_utc_time,timestamp(notAdjustedToUTC)"`
			}),
			print: `message {
	optional int64 int_time (TIMESTAMP(isAdjustedToUTC=true,unit=MILLIS));
	optional int64 time_time (TIMESTAMP(isAdjustedToUTC=true,unit=MILLIS));
	optional int64 not_utc_time (TIMESTAMP(isAdjustedToUTC=false,unit=MILLIS));
}`,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			schema := parquet.SchemaOf(test.value)

			if s := schema.String(); s != test.print {
				t.Errorf("\nexpected:\n\n%s\n\nfound:\n\n%s\n", test.print, s)
			}
		})
	}
}
