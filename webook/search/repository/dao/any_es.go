package dao

import (
	"context"
	"github.com/olivere/elastic/v7"
)

type AnyESDAO struct {
	client *elastic.Client
}

func NewAnyESDAO(client *elastic.Client) AnyDAO {
	return &AnyESDAO{client}
}

func (r *AnyESDAO) Input(ctx context.Context, index, docID, data string) error {
	_, err := r.client.Index().Index(index).Id(docID).BodyJson(data).Do(ctx)
	return err
}
