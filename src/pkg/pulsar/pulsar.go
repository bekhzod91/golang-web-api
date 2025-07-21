package pulsar

import (
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsaradmin"
	"github.com/apache/pulsar-client-go/pulsaradmin/pkg/utils"
	"sync"
)

type ConsumerArg struct {
	Admin     pulsaradmin.Client
	Client    pulsar.Client
	Tenant    string
	Namespace string
	Topic     string
}

type Topics map[string]func(message pulsar.ConsumerMessage)

func NewConsumer(tenants []string, topics Topics) error {
	pulsarAdmin, err := pulsaradmin.NewClient(&pulsaradmin.Config{
		WebServiceURL: "http://localhost:8080",
	})

	if err != nil {
		panic(fmt.Sprintf("fail connect pulsar admin %s", err.Error()))
	}

	pulsarClient, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://localhost:6650",
	})

	if err != nil {
		panic(fmt.Sprintf("fail connect pulsar %s", err.Error()))
	}

	defer pulsarClient.Close()

	err = createTenants(pulsarAdmin, tenants)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, tenant := range tenants {
		namespace := "default"
		_ = pulsarAdmin.Namespaces().CreateNamespace(fmt.Sprintf("%s/%s", tenant, namespace))

		for topic := range topics {
			consumerArg := ConsumerArg{
				Admin:     pulsarAdmin,
				Client:    pulsarClient,
				Namespace: namespace,
				Tenant:    tenant,
				Topic:     topic,
			}

			fn := topics[topic]

			wg.Add(1)
			go func() {
				defer wg.Done()
				handler(&consumerArg, fn)
			}()
		}
	}
	wg.Wait()

	//tenants, err := clientAdmin.Namespaces().N()
	//fmt.Println("tenants: ", tenants)
	//
	//if err != nil {
	//	fmt.Println("fail create pulsar client", err.Error())
	//	return err
	//}
	//
	//channel := make(chan pulsar.ConsumerMessage, 2)
	//topics := []string{
	//	"persistent://tenant1/default/create_order",
	//	"persistent://tenant2/default/create_order",
	//}
	//consumer, err := client.Subscribe(pulsar.ConsumerOptions{
	//	Topics:           topics,
	//	SubscriptionName: "create_order",
	//	MessageChannel:   channel,
	//})
	//if err != nil {
	//	fmt.Println("fail create pulsar pulsar consumer", err.Error())
	//	return err
	//}
	//
	//defer consumer.Close()
	//
	//for cm := range channel {
	//	time.Sleep(time.Duration(100) * time.Millisecond)
	//	consumer := cm.Consumer
	//	msg := cm.Message
	//	fmt.Printf("Consumer %s received a message, msgId: %v, content: '%s'\n",
	//		consumer.Name(), msg.ID(), string(msg.Payload()))
	//
	//	_ = consumer.Ack(msg)
	//}

	return nil
}

func createTenants(pulsarAdmin pulsaradmin.Client, tenants []string) error {
	pulsarTenants, err := pulsarAdmin.Tenants().List()
	if err != nil {
		return fmt.Errorf("failed to get tenant list %w", err)
	}

	tenantsByName := toMap(pulsarTenants)
	for _, tenant := range tenants {
		if t, ok := tenantsByName[tenant]; ok {
			fmt.Printf("tenant %s already exists, creation skipped!\n", t)
			continue
		}

		tenantData := utils.TenantData{
			Name:            tenant,
			AllowedClusters: []string{"standalone"},
		}
		err := pulsarAdmin.Tenants().Create(tenantData)
		if err != nil {
			return fmt.Errorf("failed to create tenant %s: %w\n", tenant, err)
		}

		pulsarAdmin.Namespaces()
	}

	return nil
}

func toMap(slice []string) map[string]string {
	result := make(map[string]string)
	for i := 0; i < len(slice); i++ {
		result[slice[i]] = slice[i]
	}
	return result
}

func handler(c *ConsumerArg, fn func(pulsar.ConsumerMessage)) {
	channel := make(chan pulsar.ConsumerMessage, 10)
	topic := fmt.Sprintf("persistent://%s/%s/%s", c.Tenant, c.Namespace, c.Topic)

	name, _ := utils.GetTopicName(topic)
	err := c.Admin.Topics().Create(*name, 0)
	//if err != nil {
	//	fmt.Printf("failed to create topic! %s", err.Error())
	//	panic(err)
	//}

	subscriptionName := fmt.Sprintf("%s_%s_%s", c.Tenant, c.Namespace, c.Topic)
	consumer, err := c.Client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subscriptionName,
		MessageChannel:   channel,
	})

	if err != nil {
		fmt.Println("fail create pulsar pulsar consumer", err.Error())
		panic(err)
	}

	defer consumer.Close()

	for cm := range channel {
		fn(cm)
		_ = consumer.Ack(cm.Message)
	}
}
