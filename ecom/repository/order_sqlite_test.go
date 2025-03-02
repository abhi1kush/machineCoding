package repository

import (
	"database/sql"
	"reflect"
	"testing"

	"ecom.com/constants"
	"ecom.com/database"
	"ecom.com/models"
	"github.com/google/uuid"
)

func TestSQLiteOrderRepository_CreateOrder(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		order *models.Order
	}
	testDb := database.ConnectDB("sqlite3", "testDb.db")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				DB: testDb,
			},
			args: args{
				order: &models.Order{
					OrderID:     "1",
					UserID:      "testUser",
					TotalAmount: 76.0,
					Status:      "Pending",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SQLiteOrderRepository{
				DB: tt.fields.DB,
			}
			if err := r.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteOrderRepository.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSQLiteOrderRepository_GetOrderByID(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		id    string
		order *models.Order
	}
	testDb := database.ConnectDB("sqlite3", "testDb.db")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name: "Basic",
			fields: fields{
				DB: testDb,
			},
			args: args{
				id: "3",
				order: &models.Order{
					OrderID:     "3",
					UserID:      "testUser",
					TotalAmount: 76.0,
					Status:      "Pending",
				},
			},
			want: &models.Order{
				OrderID:     "3",
				UserID:      "testUser",
				TotalAmount: 76.0,
				Status:      "Pending",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SQLiteOrderRepository{
				DB: tt.fields.DB,
			}
			if err := r.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteOrderRepository.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := r.GetOrderByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteOrderRepository.GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLiteOrderRepository.GetOrderByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLiteOrderRepository_UpdateOrderStatus(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		orderId string
		status  string
		order   *models.Order
	}
	newOrderId := uuid.NewString()
	testDb := database.ConnectDB("sqlite3", "testDb.db")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Order
		wantErr bool
	}{
		{
			name: "Basic",
			fields: fields{
				DB: testDb,
			},
			args: args{
				orderId: newOrderId,
				status:  string(constants.COMPELETED),
				order: &models.Order{
					OrderID:     newOrderId,
					UserID:      "testUser",
					TotalAmount: 77.0,
					Status:      "Pending",
				},
			},
			want: &models.Order{
				OrderID:     newOrderId,
				UserID:      "testUser",
				TotalAmount: 77.0,
				Status:      "Pending",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &SQLiteOrderRepository{
				DB: tt.fields.DB,
			}
			if err := r.CreateOrder(tt.args.order); (err != nil) != tt.wantErr {
				t.Errorf("SQLiteOrderRepository.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := r.GetOrderByID(tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLiteOrderRepository.GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("1 SQLiteOrderRepository.GetOrderByID() = %v, want %v", got, tt.want)
			}
			if err := r.UpdateOrderStatus(tt.args.orderId, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("2 SQLiteOrderRepository.UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotUpdated, err := r.GetOrderByID(tt.args.orderId)
			if (err != nil) != tt.wantErr {
				t.Errorf("3 SQLiteOrderRepository.GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUpdated.Status, tt.args.status) {
				t.Errorf("4 SQLiteOrderRepository.GetOrderByID() = %v, want %v", gotUpdated.Status, tt.args.status)
			}

		})
	}
}
